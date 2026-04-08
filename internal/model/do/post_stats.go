// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PostStats is the golang structure of table post_stats for DAO operations like Where/Data.
type PostStats struct {
	g.Meta       `orm:"table:post_stats, do:true"`
	PostId       any         //
	ViewCount    any         //
	LikeCount    any         //
	CommentCount any         //
	ShareCount   any         //
	UpdatedAt    *gtime.Time //
}
