// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// VerificationCodes is the golang structure for table verification_codes.
type VerificationCodes struct {
	Id        int64       `json:"id"        orm:"id"         description:""`
	Target    string      `json:"target"    orm:"target"     description:""`
	Code      string      `json:"code"      orm:"code"       description:""`
	Type      string      `json:"type"      orm:"type"       description:""`
	ExpiresAt *gtime.Time `json:"expiresAt" orm:"expires_at" description:""`
	UsedAt    *gtime.Time `json:"usedAt"    orm:"used_at"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
}
