// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// UserBookmarks is the golang structure for table user_bookmarks.
type UserBookmarks struct {
	UserId    int         `json:"userId"    orm:"user_id"    description:""` //
	PostId    int         `json:"postId"    orm:"post_id"    description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
