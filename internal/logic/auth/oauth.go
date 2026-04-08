package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	"github.com/nuxtblog/nuxtblog/internal/oauth"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gtime"
)

var oauthStateCache = gcache.New()

// OAuthProviders returns the names of all enabled OAuth providers.
func (s *sAuth) OAuthProviders(ctx context.Context) []string {
	return oauth.Enabled(ctx)
}

// OAuthRedirect generates a state token, caches it, and returns the provider's auth URL.
// redirectAfter is the frontend path to navigate to after successful login (e.g. "/").
func (s *sAuth) OAuthRedirect(ctx context.Context, providerName, redirectAfter string) (string, error) {
	if !oauth.IsEnabled(ctx, providerName) {
		return "", fmt.Errorf("%s", g.I18n().Tf(ctx, "oauth.provider_not_enabled", providerName))
	}
	p, ok := oauth.GetProvider(ctx, providerName)
	if !ok {
		return "", fmt.Errorf("%s", g.I18n().Tf(ctx, "oauth.unknown_provider", providerName))
	}
	state, err := randomHex(16)
	if err != nil {
		return "", err
	}
	if err := oauthStateCache.Set(ctx, "oauth:"+state, redirectAfter, 10*time.Minute); err != nil {
		return "", err
	}
	return p.AuthURL(state), nil
}

// OAuthCallback validates state, exchanges the code, and issues JWT tokens.
// Returns (accessToken, refreshToken, expiresIn, redirectAfter, error).
func (s *sAuth) OAuthCallback(ctx context.Context, providerName, code, state string) (string, string, int64, string, error) {
	// Validate and consume state (one-time use, CSRF protection)
	val, err := oauthStateCache.Remove(ctx, "oauth:"+state)
	if err != nil || val == nil {
		return "", "", 0, "", errors.New(g.I18n().T(ctx, "oauth.invalid_state"))
	}
	redirectAfter := val.String()

	p, ok := oauth.GetProvider(ctx, providerName)
	if !ok {
		return "", "", 0, "", fmt.Errorf("%s", g.I18n().Tf(ctx, "oauth.unknown_provider", providerName))
	}

	// Exchange code for user info
	info, err := p.Exchange(ctx, code)
	if err != nil {
		return "", "", 0, "", fmt.Errorf("%s: %w", g.I18n().T(ctx, "oauth.exchange_failed"), err)
	}

	// Find or create user
	user, err := findOrCreateOAuthUser(ctx, providerName, info)
	if err != nil {
		return "", "", 0, "", err
	}

	// Update last_login_at
	_, _ = dao.Users.Ctx(ctx).Where("id", user.Id).
		Data(g.Map{"last_login_at": gtime.Now()}).Update()

	access, refresh, expiresIn, err := issueTokens(ctx, user)
	if err != nil {
		return "", "", 0, "", err
	}
	saveRefreshToken(ctx, int64(user.Id), refresh)
	return access, refresh, expiresIn, redirectAfter, nil
}

// findOrCreateOAuthUser looks up user_oauth, links to an existing account by email,
// or creates a brand-new user.
func findOrCreateOAuthUser(ctx context.Context, providerName string, info *oauth.UserInfo) (*entity.Users, error) {
	// 1. Check if OAuth account is already linked to a user
	type oauthLinkRow struct {
		UserID int64 `orm:"user_id"`
	}
	var link oauthLinkRow
	_ = dao.UserOauth.Ctx(ctx).
		Where("provider", providerName).
		Where("provider_id", info.ProviderID).
		Scan(&link)
	if link.UserID > 0 {
		return loadEntityUser(ctx, link.UserID)
	}

	// 2. If a real email was provided, check if a user already owns it
	var userID int64
	if info.Email != "" && !strings.HasSuffix(info.Email, "@oauth.placeholder") {
		type existingRow struct {
			ID int64 `orm:"id"`
		}
		var existing existingRow
		_ = dao.Users.Ctx(ctx).WhereNull("deleted_at").
			Where("email", info.Email).Fields("id").Scan(&existing)
		userID = existing.ID
	}

	// 3. No existing user → create one
	if userID == 0 {
		username := generateUsername(ctx, providerName, info)
		now := gtime.Now()
		id, err := dao.Users.Ctx(ctx).Data(g.Map{
			"username":      username,
			"email":         info.Email,
			"password_hash": "", // OAuth-only; bcrypt never matches empty string
			"display_name":  info.Name,
			"role":          1, // subscriber
			"status":        1, // active
			"locale":        "zh-CN",
			"bio":           "",
			"created_at":    now,
			"updated_at":    now,
		}).InsertAndGetId()
		if err != nil {
			return nil, fmt.Errorf("%s: %w", g.I18n().T(ctx, "oauth.create_user_failed"), err)
		}
		userID = id
	}

	// 4. Link OAuth account → user
	now := gtime.Now()
	_, _ = dao.UserOauth.Ctx(ctx).Data(g.Map{
		"id":          idgen.New(),
		"user_id":     userID,
		"provider":    providerName,
		"provider_id": info.ProviderID,
		"created_at":  now,
	}).Insert()

	return loadEntityUser(ctx, userID)
}

func loadEntityUser(ctx context.Context, id int64) (*entity.Users, error) {
	var u entity.Users
	err := dao.Users.Ctx(ctx).Where("id", id).WhereNull("deleted_at").Scan(&u)
	if err != nil || u.Id == 0 {
		return nil, fmt.Errorf("%s", g.I18n().Tf(ctx, "oauth.user_not_found", id))
	}
	return &u, nil
}

// generateUsername creates a unique username derived from the provider and user info.
func generateUsername(ctx context.Context, providerName string, info *oauth.UserInfo) string {
	base := providerName + "_" + sanitizeUsername(info.Name)
	if base == providerName+"_" {
		base = providerName + "_" + sanitizeUsername(info.ProviderID)
	}
	if len(base) > 28 {
		base = base[:28]
	}
	candidate := base
	for i := 2; i <= 99; i++ {
		cnt, _ := dao.Users.Ctx(ctx).WhereNull("deleted_at").
			Where("username", candidate).Count()
		if cnt == 0 {
			return candidate
		}
		candidate = fmt.Sprintf("%s_%d", base, i)
	}
	return candidate
}

func sanitizeUsername(name string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '_' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
