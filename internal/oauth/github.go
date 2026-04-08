package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/nuxtblog/nuxtblog/internal/consts"
)

type githubProvider struct{}

func init() { Register(&githubProvider{}) }

func (p *githubProvider) Name() string { return consts.OAuthProviderGitHub }

func (p *githubProvider) AuthURL(state string) string {
	cfg := GetConfig(context.Background(), consts.OAuthProviderGitHub)
	q := url.Values{}
	q.Set("client_id", cfg.ClientId)
	q.Set("redirect_uri", cfg.CallbackUrl)
	q.Set("scope", "user:email")
	q.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + q.Encode()
}

func (p *githubProvider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	cfg := GetConfig(ctx, consts.OAuthProviderGitHub)

	// Exchange code for access token
	resp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {cfg.ClientId},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"redirect_uri":  {cfg.CallbackUrl},
	})
	if err != nil {
		return nil, fmt.Errorf("github token exchange: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	vals, _ := url.ParseQuery(string(body))
	accessToken := vals.Get("access_token")
	if accessToken == "" {
		return nil, fmt.Errorf("github: no access_token in response")
	}

	// Fetch user profile
	profile, err := githubRequest(accessToken, "https://api.github.com/user")
	if err != nil {
		return nil, fmt.Errorf("github user profile: %w", err)
	}
	id, _ := profile["id"].(float64)
	login, _ := profile["login"].(string)
	name, _ := profile["name"].(string)
	email, _ := profile["email"].(string)
	avatar, _ := profile["avatar_url"].(string)
	if name == "" {
		name = login
	}

	// Fetch primary verified email if not in profile
	if email == "" {
		email = githubPrimaryEmail(accessToken)
	}

	return &UserInfo{
		ProviderID: fmt.Sprintf("%d", int64(id)),
		Email:      email,
		Name:       name,
		Avatar:     avatar,
	}, nil
}

func githubRequest(token, apiURL string) (map[string]interface{}, error) {
	req, _ := http.NewRequest("GET", apiURL, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	return result, json.NewDecoder(resp.Body).Decode(&result)
}

func githubPrimaryEmail(token string) string {
	req, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var emails []struct {
		Email    string `json:"email"`
		Primary  bool   `json:"primary"`
		Verified bool   `json:"verified"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return ""
	}
	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email
		}
	}
	if len(emails) > 0 {
		return emails[0].Email
	}
	return ""
}
