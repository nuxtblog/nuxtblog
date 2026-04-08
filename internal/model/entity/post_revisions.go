// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// PostRevisions is the golang structure for table post_revisions.
type PostRevisions struct {
	Id        int         `json:"id"        orm:"id"         description:""` //
	PostId    int         `json:"postId"    orm:"post_id"    description:""` //
	AuthorId  int         `json:"authorId"  orm:"author_id"  description:""` //
	Title     string      `json:"title"     orm:"title"      description:""` //
	Content   string      `json:"content"   orm:"content"    description:""` //
	RevNote   string      `json:"revNote"   orm:"rev_note"   description:""` //
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""` //
}
