// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserLikes is the golang structure of table user_likes for DAO operations like Where/Data.
type UserLikes struct {
	g.Meta    `orm:"table:user_likes, do:true"`
	UserId    any         //
	PostId    any         //
	CreatedAt *gtime.Time //
}
