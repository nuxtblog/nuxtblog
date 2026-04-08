// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PostRevisions is the golang structure of table post_revisions for DAO operations like Where/Data.
type PostRevisions struct {
	g.Meta    `orm:"table:post_revisions, do:true"`
	Id        any         //
	PostId    any         //
	AuthorId  any         //
	Title     any         //
	Content   any         //
	RevNote   any         //
	CreatedAt *gtime.Time //
}
