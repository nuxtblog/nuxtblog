// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Notifications is the golang structure for table notifications.
type Notifications struct {
	Id          int         `json:"id"          orm:"id"           description:""` //
	UserId      int         `json:"userId"      orm:"user_id"      description:""` //
	Type        string      `json:"type"        orm:"type"         description:""` //
	SubType     string      `json:"subType"     orm:"sub_type"     description:""` //
	ActorId     int         `json:"actorId"     orm:"actor_id"     description:""` //
	ActorName   string      `json:"actorName"   orm:"actor_name"   description:""` //
	ActorAvatar string      `json:"actorAvatar" orm:"actor_avatar" description:""` //
	ObjectType  string      `json:"objectType"  orm:"object_type"  description:""` //
	ObjectId    int         `json:"objectId"    orm:"object_id"    description:""` //
	ObjectTitle string      `json:"objectTitle" orm:"object_title" description:""` //
	ObjectLink  string      `json:"objectLink"  orm:"object_link"  description:""` //
	Content     string      `json:"content"     orm:"content"      description:""` //
	IsRead      int         `json:"isRead"      orm:"is_read"      description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""` //
}
