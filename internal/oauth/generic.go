package oauth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// GenericConfig holds full OAuth2 configuration for a dynamically-added provider.
// Extends ProviderConfig (Enabled/ClientId/ClientSecret/CallbackUrl) with endpoint URLs
// and field mappings. Stored in options table under key "oauth_{slug}".
type GenericConfig struct {
	// Basic credentials (shared with builtin providers)
	Enabled      bool   `json:"enabled"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	CallbackUrl  string `json:"callbackUrl"`

	// Generic provider identity
	Slug  string `json:"slug"`
	Label string `json:"label"`
	Icon  string `json:"icon"`

	// OAuth2 endpoints
	AuthUrl     string   `json:"authUrl"`
	TokenUrl    string   `json:"tokenUrl"`
	UserInfoUrl string   `json:"userInfoUrl"`
	Scopes      []string `json:"scopes"`

	// User info field mappings (JSON field names in the userinfo response)
	Fields GenericFieldMap `json:"fields"`
}

// GenericFieldMap maps standardized UserInfo fields to the provider's actual JSON keys.
type GenericFieldMap struct {
	ID     string `json:"id"`     // e.g. "id", "sub"
	Email  string `json:"email"`  // e.g. "email"
	Name   string `json:"name"`   // e.g. "name", "username", "login"
	Avatar string `json:"avatar"` // e.g. "avatar_url", "picture"
}

// IsGeneric reports whether the config describes a generic (frontend-defined) provider.
// A generic provider must have authUrl set.
func (c *GenericConfig) IsGeneric() bool {
	return c.AuthUrl != ""
}

// genericProvider implements Provider using dynamic DB config.
type genericProvider struct {
	slug string
}

func (p *genericProvider) Name() string { return p.slug }

func (p *genericProvider) AuthURL(state string) string {
	cfg := loadGenericConfig(context.Background(), p.slug)
	if cfg == nil {
		return ""
	}
	q := url.Values{}
	q.Set("client_id", cfg.ClientId)
	q.Set("redirect_uri", cfg.CallbackUrl)
	q.Set("response_type", "code")
	q.Set("state", state)
	if len(cfg.Scopes) > 0 {
		q.Set("scope", strings.Join(cfg.Scopes, " "))
	}
	return cfg.AuthUrl + "?" + q.Encode()
}

func (p *genericProvider) Exchange(ctx context.Context, code string) (*UserInfo, error) {
	cfg := loadGenericConfig(ctx, p.slug)
	if cfg == nil {
		return nil, fmt.Errorf("%s: provider config not found", p.slug)
	}

	// Exchange code for access token
	resp, err := http.PostForm(cfg.TokenUrl, url.Values{
		"client_id":     {cfg.ClientId},
		"client_secret": {cfg.ClientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"redirect_uri":  {cfg.CallbackUrl},
	})
	if err != nil {
		return nil, fmt.Errorf("%s token exchange: %w", p.slug, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var tokenResp map[string]interface{}
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("%s token parse: %w", p.slug, err)
	}
	accessToken, _ := tokenResp["access_token"].(string)
	if accessToken == "" {
		errMsg, _ := tokenResp["error_description"].(string)
		if errMsg == "" {
			errMsg, _ = tokenResp["error"].(string)
		}
		return nil, fmt.Errorf("%s: no access_token — %s", p.slug, errMsg)
	}

	// Fetch user info
	req, _ := http.NewRequest("GET", cfg.UserInfoUrl, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	uresp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%s user info: %w", p.slug, err)
	}
	defer uresp.Body.Close()

	var profile map[string]interface{}
	if err := json.NewDecoder(uresp.Body).Decode(&profile); err != nil {
		return nil, fmt.Errorf("%s user parse: %w", p.slug, err)
	}

	fm := cfg.Fields
	id := fmt.Sprintf("%v", profile[fm.ID])
	email, _ := profile[fm.Email].(string)
	name, _ := profile[fm.Name].(string)
	avatar, _ := profile[fm.Avatar].(string)

	if email == "" {
		email = fmt.Sprintf("%s_%s@oauth.placeholder", p.slug, id)
	}
	if name == "" {
		name = id
	}

	return &UserInfo{
		ProviderID: id,
		Email:      email,
		Name:       name,
		Avatar:     avatar,
	}, nil
}

// loadGenericConfig reads a GenericConfig from the options table.
// Returns nil if not found or not a generic provider (no authUrl).
func loadGenericConfig(ctx context.Context, slug string) *GenericConfig {
	type optRow struct {
		Value string `orm:"value"`
	}
	var row optRow
	_ = g.DB().Ctx(ctx).Model("options").
		Where("key", "oauth_"+slug).
		Fields("value").
		Scan(&row)
	if row.Value == "" || row.Value == "null" {
		return nil
	}
	var cfg GenericConfig
	if err := json.Unmarshal([]byte(row.Value), &cfg); err != nil {
		return nil
	}
	if !cfg.IsGeneric() {
		return nil
	}
	cfg.Slug = slug
	return &cfg
}

// GenericProviderSlugs returns the list of custom provider slugs stored in oauth_providers option.
func GenericProviderSlugs(ctx context.Context) []string {
	type optRow struct {
		Value string `orm:"value"`
	}
	var row optRow
	_ = g.DB().Ctx(ctx).Model("options").
		Where("key", "oauth_providers").
		Fields("value").
		Scan(&row)
	if row.Value == "" {
		return nil
	}
	var slugs []string
	_ = json.Unmarshal([]byte(row.Value), &slugs)
	return slugs
}
