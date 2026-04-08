// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Options is the golang structure of table options for DAO operations like Where/Data.
type Options struct {
	g.Meta    `orm:"table:options, do:true"`
	Key       any         //
	Value     any         //
	Autoload  any         //
	UpdatedAt *gtime.Time //
}
