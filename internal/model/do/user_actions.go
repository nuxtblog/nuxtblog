// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// UserActions is the golang structure of table user_actions for DAO operations like Where/Data.
type UserActions struct {
	g.Meta     `orm:"table:user_actions, do:true"`
	Id         any         //
	UserId     any         //
	Action     any         //
	ObjectType any         //
	ObjectId   any         //
	Extra      any         //
	CreatedAt  *gtime.Time //
}
