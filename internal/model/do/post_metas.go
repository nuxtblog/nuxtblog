// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// PostMetas is the golang structure of table post_metas for DAO operations like Where/Data.
type PostMetas struct {
	g.Meta    `orm:"table:post_metas, do:true"`
	Id        any         //
	PostId    any         //
	MetaKey   any         //
	MetaValue any         //
	CreatedAt *gtime.Time //
	UpdatedAt *gtime.Time //
}
