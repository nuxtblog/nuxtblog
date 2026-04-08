// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PostStats is the golang structure for table post_stats.
type PostStats struct {
	PostId       int         `json:"postId"       orm:"post_id"       description:""` //
	ViewCount    int         `json:"viewCount"    orm:"view_count"    description:""` //
	LikeCount    int         `json:"likeCount"    orm:"like_count"    description:""` //
	CommentCount int         `json:"commentCount" orm:"comment_count" description:""` //
	ShareCount   int         `json:"shareCount"   orm:"share_count"   description:""` //
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:""` //
}
