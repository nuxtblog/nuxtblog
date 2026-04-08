// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Comments is the golang structure of table comments for DAO operations like Where/Data.
type Comments struct {
	g.Meta      `orm:"table:comments, do:true"`
	Id          any         //
	ObjectId    any         //
	ObjectType  any         //
	ParentId    any         //
	UserId      any         //
	AuthorName  any         //
	AuthorEmail any         //
	Content     any         //
	Status      any         //
	Ip          any         //
	UserAgent   any         //
	CreatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
