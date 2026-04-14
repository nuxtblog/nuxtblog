package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/tokenstore"
	"github.com/nuxtblog/nuxtblog/internal/util/password"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/golang-jwt/jwt/v5"
)

type sAuth struct{}

func New() service.IAuth { return &sAuth{} }

func init() {
	service.RegisterAuth(New())
}

// ── helpers ──────────────────────────────────────────────────────────────────

func getSecret(ctx context.Context) []byte {
	val, _ := g.Cfg().Get(ctx, "auth.jwtSecret")
	s := val.String()
	if s == "" {
		s = "change-me-in-production"
	}
	return []byte(s)
}

func getAccessExpiry(ctx context.Context) int64 {
	val, _ := g.Cfg().Get(ctx, "auth.jwtAccessExpiry")
	if n := val.Int64(); n > 0 {
		return n
	}
	return 86400
}

func getRefreshExpiry(ctx context.Context) int64 {
	val, _ := g.Cfg().Get(ctx, "auth.jwtRefreshExpiry")
	if n := val.Int64(); n > 0 {
		return n
	}
	return 2592000
}

// hashToken returns the SHA-256 hex digest of a raw token string.
// We never store raw tokens — only their hashes.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

type jwtClaims struct {
	UserID int64  `json:"user_id"`
	Role   int    `json:"role"`
	Type   string `json:"type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func issueTokens(ctx context.Context, user *entity.Users) (access, refresh string, expiresIn int64, err error) {
	secret := getSecret(ctx)
	accessExpiry := getAccessExpiry(ctx)
	refreshExpiry := getRefreshExpiry(ctx)
	now := time.Now()

	// access token
	accessClaims := jwtClaims{
		UserID: int64(user.Id),
		Role:   user.Role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.Id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(accessExpiry) * time.Second)),
		},
	}
	access, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return
	}

	// refresh token
	refreshClaims := jwtClaims{
		UserID: int64(user.Id),
		Role:   user.Role,
		Type:   "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.Id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(refreshExpiry) * time.Second)),
		},
	}
	refresh, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secret)
	expiresIn = accessExpiry
	return
}

// saveRefreshToken persists the refresh token hash to the configured store.
func saveRefreshToken(ctx context.Context, userID int64, refresh string) {
	expiry := getRefreshExpiry(ctx)
	expiresAt := time.Now().Add(time.Duration(expiry) * time.Second)
	if err := tokenstore.Default(ctx).Save(ctx, userID, hashToken(refresh), expiresAt); err != nil {
		g.Log().Warningf(ctx, "[auth] tokenstore.Save error: %v", err)
	}
}

func toUserItem(ctx context.Context, u *entity.Users) v1.AuthUserItem {
	avatarId := int64(u.AvatarId)
	item := v1.AuthUserItem{
		Id:          int64(u.Id),
		Username:    u.Username,
		Email:       u.Email,
		DisplayName: u.DisplayName,
		Bio:         u.Bio,
		Role:        u.Role,
		Status:      u.Status,
		Locale:      u.Locale,
		CreatedAt:   u.CreatedAt,
		HasPassword: u.PasswordHash != "",
	}
	if u.AvatarId > 0 {
		item.AvatarId = &avatarId
		type mediaRow struct {
			CdnUrl string `orm:"cdn_url"`
		}
		var m mediaRow
		_ = dao.Medias.Ctx(ctx).Fields("cdn_url").Where("id", u.AvatarId).Scan(&m)
		if m.CdnUrl != "" {
			item.Avatar = &m.CdnUrl
		}
	}
	return item
}

func parseToken(ctx context.Context, tokenStr string) (*jwtClaims, error) {
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return getSecret(ctx), nil
	})
	return claims, err
}

// ── Login ────────────────────────────────────────────────────────────────────

func (s *sAuth) Login(ctx context.Context, req *v1.AuthLoginReq) (*v1.AuthLoginRes, error) {
	var user entity.Users
	err := dao.Users.Ctx(ctx).
		WhereNull("deleted_at").
		Where("status", 1).
		Where("(username = ? OR email = ?)", req.Login, req.Login).
		Scan(&user)
	if err != nil || user.Id == 0 {
		return nil, errors.New(g.I18n().T(ctx, "auth.invalid_credentials"))
	}

	if !password.Verify(req.Password, user.PasswordHash) {
		return nil, errors.New(g.I18n().T(ctx, "auth.invalid_credentials"))
	}

	// Plugin filter — allows plugins to block login (e.g. IP ban, 2FA)
	if _, ferr := eng.Filter(ctx, eng.FilterUserLogin, map[string]any{
		"user_id":  user.Id,
		"username": user.Username,
		"email":    user.Email,
		"role":     user.Role,
	}); ferr != nil {
		return nil, ferr
	}

	// Update last_login_at
	_, _ = dao.Users.Ctx(ctx).Where("id", user.Id).Data(g.Map{"last_login_at": gtime.Now()}).Update()

	_ = event.Emit(ctx, event.UserLoggedIn, payload.UserLoggedIn{
		UserID:   int64(user.Id),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	})

	access, refresh, expiresIn, err := issueTokens(ctx, &user)
	if err != nil {
		return nil, err
	}
	saveRefreshToken(ctx, int64(user.Id), refresh)
	return &v1.AuthLoginRes{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         toUserItem(ctx, &user),
	}, nil
}

// ── Register ─────────────────────────────────────────────────────────────────

func (s *sAuth) Register(ctx context.Context, req *v1.AuthRegisterReq) (*v1.AuthRegisterRes, error) {
	// Check verification code if required
	mode := service.Verifycode().GetRegisterVerifyMode(ctx)
	if mode != "none" && mode != "" {
		if req.Code == "" {
			return nil, errors.New(g.I18n().T(ctx, "auth.verify_code_required"))
		}
		target := req.Email
		if mode == "sms" {
			target = req.Username // SMS uses phone number passed as username
		}
		if err := service.Verifycode().CheckCode(ctx, target, mode, req.Code); err != nil {
			return nil, err
		}
	}

	// Check uniqueness
	count, _ := dao.Users.Ctx(ctx).
		WhereNull("deleted_at").
		Where("(username = ? OR email = ?)", req.Username, req.Email).
		Count()
	if count > 0 {
		return nil, errors.New(g.I18n().T(ctx, "auth.username_email_exists"))
	}

	hash, err := password.Hash(req.Password)
	if err != nil {
		return nil, err
	}

	displayName := req.DisplayName
	if displayName == "" {
		displayName = req.Username
	}

	// Run plugin filter:user.register — allows plugins to modify or reject registration data
	username := req.Username
	email := req.Email
	if filtered, ferr := eng.Filter(ctx, eng.FilterUserRegister, map[string]any{
		"username": username,
		"email":    email,
		"display_name": displayName,
	}); ferr != nil {
		return nil, ferr
	} else {
		if v, ok := filtered["username"].(string); ok && v != "" {
			username = v
		}
		if v, ok := filtered["email"].(string); ok && v != "" {
			email = v
		}
		if v, ok := filtered["display_name"].(string); ok && v != "" {
			displayName = v
		}
	}

	now := gtime.Now()
	result, err := dao.Users.Ctx(ctx).Data(g.Map{
		"username":      username,
		"email":         email,
		"password_hash": hash,
		"display_name":  displayName,
		"role":          1, // subscriber
		"status":        1, // active
		"locale":        "zh-CN",
		"bio":           "",
		"created_at":    now,
		"updated_at":    now,
	}).InsertAndGetId()
	if err != nil {
		return nil, uniqueUserConstraintError(err)
	}

	_ = event.Emit(ctx, event.UserRegistered, payload.UserRegistered{
		UserID:   result,
		Username: username,
		Email:    email,
	})

	var user entity.Users
	_ = dao.Users.Ctx(ctx).Where("id", result).Scan(&user)

	access, refresh, expiresIn, err := issueTokens(ctx, &user)
	if err != nil {
		return nil, err
	}
	saveRefreshToken(ctx, int64(user.Id), refresh)
	return &v1.AuthRegisterRes{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		User:         toUserItem(ctx, &user),
	}, nil
}

// ── Logout ───────────────────────────────────────────────────────────────────

func (s *sAuth) Logout(ctx context.Context, refreshToken string) error {
	ts := tokenstore.Default(ctx)
	var userID int64
	r := ghttp.RequestFromCtx(ctx)
	if r != nil {
		userID = r.GetCtxVar("user_id").Int64()
	}
	var err error
	if refreshToken != "" {
		err = ts.Delete(ctx, hashToken(refreshToken))
	} else if userID > 0 {
		err = ts.DeleteByUser(ctx, userID)
	}
	_ = event.Emit(ctx, event.UserLoggedOut, payload.UserLoggedOut{UserID: userID})
	return err
}

// ── Me ───────────────────────────────────────────────────────────────────────

func (s *sAuth) Me(ctx context.Context) (*v1.AuthUserItem, error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return nil, gerror.NewCode(gcode.New(401, "", nil), g.I18n().T(ctx, "error.unauthorized"))
	}
	uid := r.GetCtxVar("user_id").Int64()
	if uid == 0 {
		return nil, gerror.NewCode(gcode.New(401, "", nil), g.I18n().T(ctx, "error.unauthorized"))
	}
	var user entity.Users
	if err := dao.Users.Ctx(ctx).Where("id", uid).WhereNull("deleted_at").Scan(&user); err != nil || user.Id == 0 {
		return nil, errors.New(g.I18n().T(ctx, "error.user_not_found"))
	}
	item := toUserItem(ctx, &user)
	return &item, nil
}

// ── Refresh ───────────────────────────────────────────────────────────────────

func (s *sAuth) Refresh(ctx context.Context, refreshToken string) (*v1.AuthRefreshRes, error) {
	claims, err := parseToken(ctx, refreshToken)
	if err != nil || claims.Type != "refresh" {
		return nil, gerror.NewCode(gcode.New(401, "", nil), g.I18n().T(ctx, "auth.invalid_refresh_token"))
	}

	// Verify the token hash exists in the store (not revoked)
	ts := tokenstore.Default(ctx)
	ok, err := ts.Exists(ctx, hashToken(refreshToken))
	if err != nil {
		return nil, fmt.Errorf("token store error: %w", err)
	}
	if !ok {
		return nil, gerror.NewCode(gcode.New(401, "", nil), g.I18n().T(ctx, "auth.refresh_token_revoked"))
	}

	var user entity.Users
	if err := dao.Users.Ctx(ctx).Where("id", claims.UserID).WhereNull("deleted_at").Scan(&user); err != nil || user.Id == 0 {
		return nil, errors.New(g.I18n().T(ctx, "error.user_not_found"))
	}

	secret := getSecret(ctx)
	expiry := getAccessExpiry(ctx)
	now := time.Now()
	accessClaims := jwtClaims{
		UserID: int64(user.Id),
		Role:   user.Role,
		Type:   "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprintf("%d", user.Id),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiry) * time.Second)),
		},
	}
	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &v1.AuthRefreshRes{AccessToken: access, ExpiresIn: expiry}, nil
}


func uniqueUserConstraintError(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	switch {
	case strings.Contains(msg, "users.email") || (strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "email")):
		return fmt.Errorf("%s", g.I18n().T(context.Background(), "auth.email_registered"))
	case strings.Contains(msg, "users.username") || (strings.Contains(msg, "UNIQUE") && strings.Contains(msg, "username")):
		return fmt.Errorf("%s", g.I18n().T(context.Background(), "auth.username_exists"))
	case strings.Contains(msg, "UNIQUE") || strings.Contains(msg, "Duplicate entry") || strings.Contains(msg, "unique"):
		return fmt.Errorf("%s", g.I18n().T(context.Background(), "auth.username_email_exists"))
	}
	return err
}
