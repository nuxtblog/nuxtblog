// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// Taxonomies is the golang structure of table taxonomies for DAO operations like Where/Data.
type Taxonomies struct {
	g.Meta      `orm:"table:taxonomies, do:true"`
	Id          any //
	TermId      any //
	Taxonomy    any //
	Description any //
	ParentId    any //
	PostCount   any //
	Extra       any //
}
