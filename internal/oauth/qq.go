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

type qqProvider struct{}

func init() { Register(&qqProvider{}) }

func (p *qqProvider) Name() string { return consts.OAuthProviderQQ }

func (p *qqProvider) AuthURL(state string) string {
	cfg := GetConfig(context.Background(), consts.OAuthProviderQQ)
	q := url.Values{}
	q.Set("response_type", "code")
	q.Set("client_id", cfg.ClientId)
	q.Set("redirect_uri", cfg.CallbackUrl)
	q.Set("state", state)
	return "https://graph.qq.com/oauth2.0/authorize?" + q.Encode()
}

func (p *qqProvider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	cfg := GetConfig(ctx, consts.OAuthProviderQQ)

	// Exchange code for access token (QQ returns URL-encoded response)
	tokenURL := "https://graph.qq.com/oauth2.0/token?" + url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {cfg.ClientId},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"redirect_uri":  {cfg.CallbackUrl},
		"fmt":           {"json"},
	}.Encode()
	resp, err := http.Get(tokenURL)
	if err != nil {
		return nil, fmt.Errorf("qq token exchange: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	// Try JSON first (when fmt=json), fall back to URL-encoded
	accessToken := ""
	var jsonResp struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &jsonResp); err == nil && jsonResp.AccessToken != "" {
		accessToken = jsonResp.AccessToken
	} else {
		vals, _ := url.ParseQuery(string(body))
		accessToken = vals.Get("access_token")
	}
	if accessToken == "" {
		return nil, fmt.Errorf("qq: no access_token in response")
	}

	// Get OpenID — QQ returns JSONP: callback( {"client_id":"...","openid":"..."} );\n
	openID, err := qqGetOpenID(accessToken)
	if err != nil {
		return nil, err
	}

	// Get user info
	infoURL := "https://graph.qq.com/user/get_user_info?" + url.Values{
		"access_token":       {accessToken},
		"oauth_consumer_key": {cfg.ClientId},
		"openid":             {openID},
	}.Encode()
	iresp, err := http.Get(infoURL)
	if err != nil {
		return nil, fmt.Errorf("qq user info: %w", err)
	}
	defer iresp.Body.Close()
	var info struct {
		Ret        int    `json:"ret"`
		Nickname   string `json:"nickname"`
		FigureURL2 string `json:"figureurl_qq_2"` // 100x100 avatar
		FigureURL1 string `json:"figureurl_qq_1"` // 40x40 avatar
	}
	if err := json.NewDecoder(iresp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("qq user parse: %w", err)
	}
	if info.Ret != 0 {
		return nil, fmt.Errorf("qq user info error: ret=%d", info.Ret)
	}
	avatar := info.FigureURL2
	if avatar == "" {
		avatar = info.FigureURL1
	}

	// QQ does not provide email; use placeholder
	return &UserInfo{
		ProviderID: openID,
		Email:      fmt.Sprintf("qq_%s@oauth.placeholder", openID),
		Name:       info.Nickname,
		Avatar:     avatar,
	}, nil
}

// qqGetOpenID strips the JSONP wrapper and extracts the openid.
func qqGetOpenID(accessToken string) (string, error) {
	resp, err := http.Get("https://graph.qq.com/oauth2.0/me?access_token=" + url.QueryEscape(accessToken) + "&unionid=1")
	if err != nil {
		return "", fmt.Errorf("qq openid: %w", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	s := strings.TrimSpace(string(body))
	// Strip callback( ... );\n wrapper
	if strings.HasPrefix(s, "callback(") {
		s = strings.TrimPrefix(s, "callback(")
		s = strings.TrimSuffix(s, ");")
		s = strings.TrimSpace(s)
	}
	var result struct {
		OpenID string `json:"openid"`
	}
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return "", fmt.Errorf("qq openid parse: %w", err)
	}
	if result.OpenID == "" {
		return "", fmt.Errorf("qq: empty openid")
	}
	return result.OpenID, nil
}
