// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Users is the golang structure for table users.
type Users struct {
	Id           int         `json:"id"           orm:"id"            description:""` //
	Username     string      `json:"username"     orm:"username"      description:""` //
	Email        string      `json:"email"        orm:"email"         description:""` //
	PasswordHash string      `json:"passwordHash" orm:"password_hash" description:""` //
	DisplayName  string      `json:"displayName"  orm:"display_name"  description:""` //
	AvatarId     int         `json:"avatarId"     orm:"avatar_id"     description:""` //
	Bio          string      `json:"bio"          orm:"bio"           description:""` //
	Role         int         `json:"role"         orm:"role"          description:""` //
	Status       int         `json:"status"       orm:"status"        description:""` //
	Locale       string      `json:"locale"       orm:"locale"        description:""` //
	LastLoginAt  *gtime.Time `json:"lastLoginAt"  orm:"last_login_at" description:""` //
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:""` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""` //
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:""` //
}
