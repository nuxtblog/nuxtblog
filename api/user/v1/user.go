package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ----------------------------------------------------------------
//  Enums
// ----------------------------------------------------------------

type UserRole int

const (
	UserRoleSubscriber UserRole = 1
	UserRoleEditor     UserRole = 2
	UserRoleAdmin      UserRole = 3
	UserRoleSuperAdmin UserRole = 4
)

type UserStatus int

const (
	UserStatusActive  UserStatus = 1
	UserStatusBanned  UserStatus = 2
	UserStatusPending UserStatus = 3
)

// ----------------------------------------------------------------
//  Shared output
// ----------------------------------------------------------------

type UserItem struct {
	Id            int64             `json:"id"`
	Username      string            `json:"username"`
	Email         string            `json:"email"`
	DisplayName   string            `json:"display_name"`
	AvatarId      *int64            `json:"avatar_id,omitempty"`
	Avatar        *string           `json:"avatar,omitempty"`
	Cover         *string           `json:"cover,omitempty"`
	Bio           string            `json:"bio"`
	Role          UserRole          `json:"role"`
	Status        UserStatus        `json:"status"`
	EmailVerified int               `json:"email_verified"`
	Locale        string            `json:"locale"`
	Metas         map[string]string `json:"metas,omitempty"`
	LastLoginAt   *gtime.Time       `json:"last_login_at,omitempty"`
	CreatedAt     *gtime.Time       `json:"created_at,omitempty"`
	UpdatedAt     *gtime.Time       `json:"updated_at,omitempty"`
}

// ----------------------------------------------------------------
//  Create
// ----------------------------------------------------------------

type UserCreateReq struct {
	g.Meta      `path:"/users" method:"post" tags:"User" summary:"Create user"`
	Username    string   `v:"required|length:3,30"        dc:"username"`
	Email       string   `v:"required|email"              dc:"email"`
	Password    string   `v:"required|length:8,64"        dc:"plain password"`
	DisplayName string   `v:"required|length:1,100"       dc:"display name"`
	Role        UserRole `v:"required|in:1,2,3,4"    dc:"1=subscriber 2=editor 3=admin 4=super_admin"`
	Locale      string   `v:"length:2,10"                 dc:"e.g. zh-CN"`
}
type UserCreateRes struct {
	Id int64 `json:"id"`
}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type UserDeleteReq struct {
	g.Meta `path:"/users/{id}" method:"delete" tags:"User" summary:"Delete user"`
	Id     int64 `v:"required|min:1" dc:"user id"`
}
type UserDeleteRes struct{}

// ----------------------------------------------------------------
//  Update
// ----------------------------------------------------------------

type UserUpdateReq struct {
	g.Meta      `path:"/users/{id}" method:"put" tags:"User" summary:"Update user"`
	Id          int64       `v:"required|min:1"          dc:"user id"`
	DisplayName *string     `v:"length:1,100"            dc:"display name"`
	Bio         *string     `v:"max-length:500"          dc:"bio"`
	AvatarId    *int64      `v:"min:1"                   dc:"avatar media id"`
	Locale      *string     `v:"length:2,10"             dc:"e.g. zh-CN"`
	Status      *UserStatus `v:"in:1,2,3"               dc:"1=active 2=banned 3=pending"`
	Role        *UserRole   `v:"in:1,2,3,4"              dc:"1=subscriber 2=editor 3=admin 4=super_admin"`
	// Profile fields (saved to user_profiles table)
	Location  *string `json:"location,omitempty"  dc:"location string"`
	Website   *string `json:"website,omitempty"   dc:"personal website URL"`
	Github    *string `json:"github,omitempty"    dc:"github username"`
	Twitter   *string `json:"twitter,omitempty"   dc:"twitter/x username"`
	Instagram *string `json:"instagram,omitempty" dc:"instagram username"`
	Linkedin  *string `json:"linkedin,omitempty"  dc:"linkedin username"`
	Youtube   *string `json:"youtube,omitempty"   dc:"youtube handle"`
	Cover     *string `json:"cover,omitempty"     dc:"cover/banner image URL"`
	CoverId   *int64  `json:"cover_id,omitempty"  dc:"cover media id for deletion tracking"`
}
type UserUpdateRes struct{}

// ----------------------------------------------------------------
//  Change password
// ----------------------------------------------------------------

type UserChangePasswordReq struct {
	g.Meta      `path:"/users/{id}/password" method:"put" tags:"User" summary:"Change password"`
	Id          int64  `v:"required|min:1"       dc:"user id"`
	OldPassword string `dc:"current password (required when changing own password)"`
	NewPassword string `v:"required|length:8,64" dc:"new password"`
}
type UserChangePasswordRes struct{}

// ----------------------------------------------------------------
//  Get one
// ----------------------------------------------------------------

type UserGetOneReq struct {
	g.Meta `path:"/users/{id}" method:"get" tags:"User" summary:"Get user by id"`
	Id     int64 `v:"required|min:1" dc:"user id"`
}
type UserGetOneRes struct {
	*UserItem `dc:"user"`
}

// ----------------------------------------------------------------
//  Get list
// ----------------------------------------------------------------

type UserGetListReq struct {
	g.Meta  `path:"/users" method:"get" tags:"User" summary:"Get user list"`
	Role    *UserRole   `v:"in:1,2,3,4"  dc:"filter by role"`
	Status  *UserStatus `v:"in:1,2,3"    dc:"filter by status"`
	Keyword *string     `                dc:"search username or email"`
	Page    int         `v:"min:1"       dc:"page number, default 1"  d:"1"`
	Size    int         `v:"between:1,100" dc:"page size, default 20" d:"20"`
}
type UserGetListRes struct {
	List  []*UserItem `json:"list"`
	Total int         `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}
