// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
)

// PostSeo is the golang structure of table post_seo for DAO operations like Where/Data.
type PostSeo struct {
	g.Meta         `orm:"table:post_seo, do:true"`
	PostId         any //
	MetaTitle      any //
	MetaDesc       any //
	OgTitle        any //
	OgImage        any //
	CanonicalUrl   any //
	Robots         any //
	StructuredData any //
}
