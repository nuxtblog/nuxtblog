package entity

import "github.com/gogf/gf/v2/os/gtime"

type MomentStats struct {
	MomentId     int         `json:"momentId"     orm:"moment_id"`
	ViewCount    int64       `json:"viewCount"    orm:"view_count"`
	LikeCount    int64       `json:"likeCount"    orm:"like_count"`
	CommentCount int64       `json:"commentCount" orm:"comment_count"`
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"`
}
