package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// AuthUserItem is a safe user representation (no password hash)
type AuthUserItem struct {
	Id          int64       `json:"id"`
	Username    string      `json:"username"`
	Email       string      `json:"email"`
	DisplayName string      `json:"display_name"`
	AvatarId    *int64      `json:"avatar_id,omitempty"`
	Avatar      *string     `json:"avatar,omitempty"`
	Bio         string      `json:"bio"`
	Role        int         `json:"role"`
	Status      int         `json:"status"`
	Locale      string      `json:"locale"`
	CreatedAt   *gtime.Time `json:"created_at"`
	HasPassword bool        `json:"has_password"` // false = OAuth-only, no password set yet
}

// ----------------------------------------------------------------
//  Login
// ----------------------------------------------------------------

type AuthLoginReq struct {
	g.Meta   `path:"/auth/login" method:"post" tags:"Auth" summary:"Login"`
	Login    string `v:"required" dc:"username or email"`
	Password string `v:"required" dc:"plain password"`
}
type AuthLoginRes struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         AuthUserItem `json:"user"`
}

// ----------------------------------------------------------------
//  Register
// ----------------------------------------------------------------

type AuthRegisterReq struct {
	g.Meta      `path:"/auth/register" method:"post" tags:"Auth" summary:"Register new user"`
	Username    string `v:"required|length:3,30"  dc:"username (3-30 chars)"`
	Email       string `v:"required|email"        dc:"email address"`
	Password    string `v:"required|length:8,64"  dc:"password (8-64 chars)"`
	DisplayName string `v:"max-length:100"        dc:"display name"`
	Code        string `dc:"verification code (required when verify mode != none)"`
}
type AuthRegisterRes struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
	User         AuthUserItem `json:"user"`
}

// ----------------------------------------------------------------
//  Logout
// ----------------------------------------------------------------

type AuthLogoutReq struct {
	g.Meta       `path:"/auth/logout" method:"post" tags:"Auth" summary:"Logout"`
	RefreshToken string `json:"refresh_token" dc:"refresh token to revoke (optional)"`
}
type AuthLogoutRes struct{}

// ----------------------------------------------------------------
//  Refresh token
// ----------------------------------------------------------------

type AuthRefreshReq struct {
	g.Meta       `path:"/auth/refresh" method:"post" tags:"Auth" summary:"Refresh access token"`
	RefreshToken string `v:"required" json:"refresh_token" dc:"refresh token"`
}
type AuthRefreshRes struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

// ----------------------------------------------------------------
//  Get current user (me)
// ----------------------------------------------------------------

type AuthMeReq struct {
	g.Meta `path:"/auth/me" method:"get" tags:"Auth" summary:"Get current logged-in user"`
}
type AuthMeRes struct {
	User AuthUserItem `json:"user"`
}

// ----------------------------------------------------------------
//  OAuth — list enabled providers
// ----------------------------------------------------------------

type AuthOAuthProvidersReq struct {
	g.Meta `path:"/auth/oauth/providers" method:"get" tags:"Auth" summary:"List enabled OAuth login providers"`
}
type AuthOAuthProvidersRes struct {
	Providers []string `json:"providers" dc:"enabled provider names, e.g. [\"github\",\"google\"]"`
}
