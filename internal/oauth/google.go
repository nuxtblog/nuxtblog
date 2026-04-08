package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nuxtblog/nuxtblog/internal/consts"
)

type googleProvider struct{}

func init() { Register(&googleProvider{}) }

func (p *googleProvider) Name() string { return consts.OAuthProviderGoogle }

func (p *googleProvider) AuthURL(state string) string {
	cfg := GetConfig(context.Background(), consts.OAuthProviderGoogle)
	q := url.Values{}
	q.Set("client_id", cfg.ClientId)
	q.Set("redirect_uri", cfg.CallbackUrl)
	q.Set("response_type", "code")
	q.Set("scope", "openid email profile")
	q.Set("state", state)
	q.Set("access_type", "online")
	return "https://accounts.google.com/o/oauth2/v2/auth?" + q.Encode()
}

func (p *googleProvider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	cfg := GetConfig(ctx, consts.OAuthProviderGoogle)

	// Exchange code for access token
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {cfg.ClientId},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {cfg.CallbackUrl},
	})
	if err != nil {
		return nil, fmt.Errorf("google token exchange: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var tokenResp struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
	}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("google token parse: %w", err)
	}
	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("google: no access_token: %s", tokenResp.Error)
	}

	// Fetch user info
	req, _ := http.NewRequest("GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	req.Header.Set("Authorization", "Bearer "+tokenResp.AccessToken)
	uresp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("google user info: %w", err)
	}
	defer uresp.Body.Close()
	var profile struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(uresp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("google user parse: %w", err)
	}

	name := profile.Name
	if name == "" && profile.Email != "" {
		name = strings.SplitN(profile.Email, "@", 2)[0]
	}
	return &UserInfo{
		ProviderID: profile.ID,
		Email:      profile.Email,
		Name:       name,
		Avatar:     profile.Picture,
	}, nil
}
