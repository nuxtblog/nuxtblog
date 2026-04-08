// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Terms is the golang structure of table terms for DAO operations like Where/Data.
type Terms struct {
	g.Meta    `orm:"table:terms, do:true"`
	Id        any         //
	Name      any         //
	Slug      any         //
	CreatedAt *gtime.Time //
}
