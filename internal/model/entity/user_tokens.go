// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// UserTokens is the golang structure for table user_tokens.
type UserTokens struct {
	Id         int64       `json:"id"         orm:"id"           description:""`
	UserId     int64       `json:"userId"     orm:"user_id"      description:""`
	Name       string      `json:"name"       orm:"name"         description:""`
	Prefix     string      `json:"prefix"     orm:"prefix"       description:""`
	TokenHash  string      `json:"tokenHash"  orm:"token_hash"   description:""`
	ExpiresAt  *gtime.Time `json:"expiresAt"  orm:"expires_at"   description:""`
	LastUsedAt *gtime.Time `json:"lastUsedAt" orm:"last_used_at" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"   description:""`
}
