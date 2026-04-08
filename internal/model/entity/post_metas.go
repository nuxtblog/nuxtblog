// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PostMetas is the golang structure for table post_metas.
type PostMetas struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	PostId    int         `json:"postId"    orm:"post_id"    description:""` //
	MetaKey   string      `json:"metaKey"   orm:"meta_key"   description:""` //
	MetaValue string      `json:"metaValue" orm:"meta_value" description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""` //
}
