// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Posts is the golang structure for table posts.
type Posts struct {
	Id            int         `json:"id"            orm:"id"              description:""` //
	PostType      int         `json:"postType"      orm:"post_type"       description:""` //
	Status        int         `json:"status"        orm:"status"          description:""` //
	Title         string      `json:"title"         orm:"title"           description:""` //
	Slug          string      `json:"slug"          orm:"slug"            description:""` //
	Content       string      `json:"content"       orm:"content"         description:""` //
	Excerpt       string      `json:"excerpt"       orm:"excerpt"         description:""` //
	AuthorId      int         `json:"authorId"      orm:"author_id"       description:""` //
	FeaturedImgId int         `json:"featuredImgId" orm:"featured_img_id" description:""` //
	CommentStatus int         `json:"commentStatus" orm:"comment_status"  description:""` //
	PasswordHash  string      `json:"passwordHash"  orm:"password_hash"   description:""` //
	Locale        string      `json:"locale"        orm:"locale"          description:""` //
	PublishedAt   *gtime.Time `json:"publishedAt"   orm:"published_at"    description:""` //
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"      description:""` //
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"      description:""` //
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"      description:""` //
}
