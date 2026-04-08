package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ----------------------------------------------------------------
//  Shared output
// ----------------------------------------------------------------

type TermItem struct {
	Id        int64       `json:"id"`
	Name      string      `json:"name"`
	Slug      string      `json:"slug"`
	CreatedAt *gtime.Time `json:"created_at"`
}

type TaxonomyItem struct {
	Id          int64           `json:"id"`
	TermId      int64           `json:"term_id"`
	Taxonomy    string          `json:"taxonomy"`
	Description string          `json:"description"`
	ParentId    *int64          `json:"parent_id"`
	PostCount   int             `json:"post_count"`
	Extra       string          `json:"extra"` // JSON
	Term        *TermItem       `json:"term,omitempty"`
	Children    []*TaxonomyItem `json:"children,omitempty"`
}

// ----------------------------------------------------------------
//  Term — Create
// ----------------------------------------------------------------

type TermCreateReq struct {
	g.Meta `path:"/terms" method:"post" tags:"Taxonomy" summary:"Create term"`
	Name   string `v:"required|length:1,100"                          dc:"term name"`
	Slug   string `v:"required|length:1,100|regex:^[a-z0-9-]+$"      dc:"URL slug"`
}
type TermCreateRes struct {
	Id int64 `json:"id"`
}

// ----------------------------------------------------------------
//  Term — Delete
// ----------------------------------------------------------------

type TermDeleteReq struct {
	g.Meta `path:"/terms/{id}" method:"delete" tags:"Taxonomy" summary:"Delete term"`
	Id     int64 `v:"required|min:1" dc:"term id"`
}
type TermDeleteRes struct{}

// ----------------------------------------------------------------
//  Term — Update
// ----------------------------------------------------------------

type TermUpdateReq struct {
	g.Meta `path:"/terms/{id}" method:"put" tags:"Taxonomy" summary:"Update term"`
	Id     int64   `v:"required|min:1"                             dc:"term id"`
	Name   *string `v:"length:1,100"                               dc:"term name"`
	Slug   *string `v:"length:1,100|regex:^[a-z0-9-]+$"           dc:"URL slug"`
}
type TermUpdateRes struct{}

// ----------------------------------------------------------------
//  Taxonomy — Create
// ----------------------------------------------------------------

type TaxonomyCreateReq struct {
	g.Meta      `path:"/taxonomies" method:"post" tags:"Taxonomy" summary:"Create taxonomy"`
	TermId      int64  `v:"required|min:1"        dc:"term id"`
	Taxonomy    string `v:"required|length:1,50"  dc:"category / tag / topic …"`
	Description string `v:"max-length:255"        dc:"description"`
	ParentId    *int64 `v:"min:1"                 dc:"parent taxonomy id"`
	Extra       string `                          dc:"JSON extra attributes"`
}
type TaxonomyCreateRes struct {
	Id int64 `json:"id"`
}

// ----------------------------------------------------------------
//  Taxonomy — Delete
// ----------------------------------------------------------------

type TaxonomyDeleteReq struct {
	g.Meta `path:"/taxonomies/{id}" method:"delete" tags:"Taxonomy" summary:"Delete taxonomy"`
	Id     int64 `v:"required|min:1" dc:"taxonomy id"`
}
type TaxonomyDeleteRes struct{}

// ----------------------------------------------------------------
//  Taxonomy — Update
// ----------------------------------------------------------------

type TaxonomyUpdateReq struct {
	g.Meta      `path:"/taxonomies/{id}" method:"put" tags:"Taxonomy" summary:"Update taxonomy"`
	Id          int64   `v:"required|min:1"       dc:"taxonomy id"`
	Description *string `v:"max-length:255"       dc:"description"`
	ParentId    *int64  `v:"min:0"                dc:"parent taxonomy id; 0 = top-level (clears parent)"`
	Extra       *string `                         dc:"JSON extra attributes"`
}
type TaxonomyUpdateRes struct{}

// ----------------------------------------------------------------
//  Taxonomy — Get tree
// ----------------------------------------------------------------

type TaxonomyGetTreeReq struct {
	g.Meta   `path:"/taxonomies/tree" method:"get" tags:"Taxonomy" summary:"Get taxonomy tree"`
	Taxonomy string `v:"required|length:1,50" dc:"e.g. category"`
}
type TaxonomyGetTreeRes struct {
	List []*TaxonomyItem `json:"list"`
}

// ----------------------------------------------------------------
//  Taxonomy — Get list (flat)
// ----------------------------------------------------------------

type TaxonomyGetListReq struct {
	g.Meta   `path:"/taxonomies" method:"get" tags:"Taxonomy" summary:"Get taxonomy list"`
	Taxonomy string `v:"required|length:1,50" dc:"e.g. category / tag"`
	Page     int    `v:"min:1"                dc:"page number" d:"1"`
	Size     int    `v:"between:1,100"        dc:"page size"   d:"50"`
}
type TaxonomyGetListRes struct {
	List  []*TaxonomyItem `json:"list"`
	Total int             `json:"total"`
}

// ----------------------------------------------------------------
//  Object-Taxonomy — bind / unbind
// ----------------------------------------------------------------

type ObjectTaxonomyBindReq struct {
	g.Meta      `path:"/object-taxonomies/bind" method:"post" tags:"Taxonomy" summary:"Bind object to taxonomies"`
	ObjectId    int64   `v:"required|min:1"        dc:"object id, e.g. post id"`
	ObjectType  string  `v:"required|in:post,product,video" dc:"object type"`
	TaxonomyIds []int64 `v:"required|min-length:1" dc:"taxonomy ids to bind"`
}
type ObjectTaxonomyBindRes struct{}

type ObjectTaxonomyUnbindReq struct {
	g.Meta      `path:"/object-taxonomies/unbind" method:"post" tags:"Taxonomy" summary:"Unbind object from taxonomies"`
	ObjectId    int64   `v:"required|min:1"        dc:"object id"`
	ObjectType  string  `v:"required|in:post,product,video" dc:"object type"`
	TaxonomyIds []int64 `v:"required|min-length:1" dc:"taxonomy ids to unbind"`
}
type ObjectTaxonomyUnbindRes struct{}

type ObjectTaxonomyGetReq struct {
	g.Meta     `path:"/object-taxonomies" method:"get" tags:"Taxonomy" summary:"Get taxonomies for an object"`
	ObjectId   int64  `v:"required|min:1"                 dc:"object id"`
	ObjectType string `v:"required|in:post,product,video" dc:"object type"`
	Taxonomy   string `v:"length:1,50"                    dc:"filter by taxonomy type"`
}
type ObjectTaxonomyGetRes struct {
	List []*TaxonomyItem `json:"list"`
}
