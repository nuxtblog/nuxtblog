// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Posts is the golang structure of table posts for DAO operations like Where/Data.
type Posts struct {
	g.Meta        `orm:"table:posts, do:true"`
	Id            any         //
	PostType      any         //
	Status        any         //
	Title         any         //
	Slug          any         //
	Content       any         //
	Excerpt       any         //
	AuthorId      any         //
	FeaturedImgId any         //
	CommentStatus any         //
	PasswordHash  any         //
	Locale        any         //
	PublishedAt   *gtime.Time //
	CreatedAt     *gtime.Time //
	UpdatedAt     *gtime.Time //
	DeletedAt     *gtime.Time //
}
