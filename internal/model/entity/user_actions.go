// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserActions is the golang structure for table user_actions.
type UserActions struct {
	Id         int         `json:"id"         orm:"id"          description:""` //
	UserId     int         `json:"userId"     orm:"user_id"     description:""` //
	Action     string      `json:"action"     orm:"action"      description:""` //
	ObjectType string      `json:"objectType" orm:"object_type" description:""` //
	ObjectId   int         `json:"objectId"   orm:"object_id"   description:""` //
	Extra      string      `json:"extra"      orm:"extra"       description:""` //
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"  description:""` //
}
