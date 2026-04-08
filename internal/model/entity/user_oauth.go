// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// UserOauth is the golang structure for table user_oauth.
type UserOauth struct {
	Id         int64       `json:"id"         orm:"id"          description:""`
	UserId     int64       `json:"userId"     orm:"user_id"     description:""`
	Provider   string      `json:"provider"   orm:"provider"    description:""`
	ProviderId string      `json:"providerId" orm:"provider_id" description:""`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""`
}
