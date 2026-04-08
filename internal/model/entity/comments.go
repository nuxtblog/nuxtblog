// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Comments is the golang structure for table comments.
type Comments struct {
	Id          int         `json:"id"          orm:"id"           description:""` //
	ObjectId    int         `json:"objectId"    orm:"object_id"    description:""` //
	ObjectType  string      `json:"objectType"  orm:"object_type"  description:""` //
	ParentId    int         `json:"parentId"    orm:"parent_id"    description:""` //
	UserId      int         `json:"userId"      orm:"user_id"      description:""` //
	AuthorName  string      `json:"authorName"  orm:"author_name"  description:""` //
	AuthorEmail string      `json:"authorEmail" orm:"author_email" description:""` //
	Content     string      `json:"content"     orm:"content"      description:""` //
	Status      int         `json:"status"      orm:"status"       description:""` //
	Ip          string      `json:"ip"          orm:"ip"           description:""` //
	UserAgent   string      `json:"userAgent"   orm:"user_agent"   description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""` //
}
