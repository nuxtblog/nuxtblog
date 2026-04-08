// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserBookmarks is the golang structure of table user_bookmarks for DAO operations like Where/Data.
type UserBookmarks struct {
	g.Meta    `orm:"table:user_bookmarks, do:true"`
	UserId    any         //
	PostId    any         //
	CreatedAt *gtime.Time //
}
