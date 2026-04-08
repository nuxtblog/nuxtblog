package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ----------------------------------------------------------------
//  Enums
// ----------------------------------------------------------------

type DocStatus int

const (
	DocStatusDraft     DocStatus = 1
	DocStatusPublished DocStatus = 2
	DocStatusArchived  DocStatus = 3
)

type CollectionStatus int

const (
	CollectionStatusDraft     CollectionStatus = 1
	CollectionStatusPublished CollectionStatus = 2
)

// ----------------------------------------------------------------
//  Shared output structs
// ----------------------------------------------------------------

type DocCollectionItem struct {
	Id          int64       `json:"id"`
	Slug        string      `json:"slug"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	CoverImgId  *int64      `json:"cover_img_id"`
	AuthorId    int64       `json:"author_id"`
	Status      int         `json:"status"`
	Locale      string      `json:"locale"`
	SortOrder   int         `json:"sort_order"`
	CreatedAt   *gtime.Time `json:"created_at"`
	UpdatedAt   *gtime.Time `json:"updated_at"`
	DocCount    int         `json:"doc_count,omitempty"`
}

type DocItem struct {
	Id            int64         `json:"id"`
	CollectionId  int64         `json:"collection_id"`
	ParentId      *int64        `json:"parent_id"`
	SortOrder     int           `json:"sort_order"`
	Status        int           `json:"status"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Excerpt       string        `json:"excerpt"`
	AuthorId      int64         `json:"author_id"`
	CommentStatus int           `json:"comment_status"`
	Locale        string        `json:"locale"`
	PublishedAt   *gtime.Time   `json:"published_at"`
	CreatedAt     *gtime.Time   `json:"created_at"`
	UpdatedAt     *gtime.Time   `json:"updated_at"`
	Stats         *DocStatsItem `json:"stats,omitempty"`
}

type DocDetailItem struct {
	DocItem
	Content string      `json:"content"`
	Seo     *DocSeoItem `json:"seo,omitempty"`
}

type DocStatsItem struct {
	ViewCount    int64 `json:"view_count"`
	LikeCount    int64 `json:"like_count"`
	CommentCount int64 `json:"comment_count"`
}

type DocSeoItem struct {
	MetaTitle      string `json:"meta_title"`
	MetaDesc       string `json:"meta_desc"`
	OgTitle        string `json:"og_title"`
	OgImage        string `json:"og_image"`
	CanonicalUrl   string `json:"canonical_url"`
	Robots         string `json:"robots"`
	StructuredData string `json:"structured_data"`
}

type DocRevisionItem struct {
	Id        int64       `json:"id"`
	DocId     int64       `json:"doc_id"`
	AuthorId  int64       `json:"author_id"`
	Title     string      `json:"title"`
	RevNote   string      `json:"rev_note"`
	CreatedAt *gtime.Time `json:"created_at"`
}

// ----------------------------------------------------------------
//  Collection CRUD
// ----------------------------------------------------------------

type CollectionCreateReq struct {
	g.Meta      `path:"/doc-collections" method:"post" tags:"Doc" summary:"Create doc collection"`
	Slug        string `v:"required|length:1,200|regex:^[a-z0-9-]+$" dc:"URL slug"`
	Title       string `v:"required|length:1,200"                     dc:"collection title"`
	Description string `v:"max-length:1000"                           dc:"description"`
	CoverImgId  *int64 `v:"min:1"                                     dc:"cover image media id"`
	Status      int    `v:"required|in:1,2"                           dc:"1=draft 2=published"`
	Locale      string `v:"length:2,10"                               dc:"e.g. zh-CN"`
	SortOrder   int    `                                               dc:"sort order"`
}
type CollectionCreateRes struct {
	Id int64 `json:"id"`
}

type CollectionUpdateReq struct {
	g.Meta      `path:"/doc-collections/{id}" method:"put" tags:"Doc" summary:"Update doc collection"`
	Id          int64   `v:"required|min:1"                            dc:"collection id"`
	Slug        *string `v:"length:1,200|regex:^[a-z0-9-]+$"          dc:"URL slug"`
	Title       *string `v:"length:1,200"                              dc:"collection title"`
	Description *string `v:"max-length:1000"                           dc:"description"`
	CoverImgId  *int64  `v:"min:1"                                     dc:"cover image media id"`
	Status      *int    `v:"in:1,2"                                    dc:"1=draft 2=published"`
	Locale      *string `v:"length:2,10"                               dc:"locale"`
	SortOrder   *int    `                                               dc:"sort order"`
}
type CollectionUpdateRes struct{}

type CollectionDeleteReq struct {
	g.Meta `path:"/doc-collections/{id}" method:"delete" tags:"Doc" summary:"Delete doc collection"`
	Id     int64 `v:"required|min:1" dc:"collection id"`
}
type CollectionDeleteRes struct{}

type CollectionGetOneReq struct {
	g.Meta `path:"/doc-collections/{id}" method:"get" tags:"Doc" summary:"Get collection by id"`
	Id     int64 `v:"required|min:1" dc:"collection id"`
}
type CollectionGetOneRes struct {
	*DocCollectionItem `dc:"collection detail"`
}

type CollectionGetListReq struct {
	g.Meta   `path:"/doc-collections" method:"get" tags:"Doc" summary:"List doc collections"`
	Status   *int    `p:"status" v:"in:1,2"           dc:"filter by status"`
	Locale   *string `p:"locale" v:"length:2,10"      dc:"filter by locale"`
	AuthorId *int64  `p:"author_id" v:"min:1"         dc:"filter by author"`
	Page     int     `p:"page" v:"min:1" d:"1"        dc:"page number"`
	PageSize int     `p:"page_size" v:"between:1,100" d:"20" dc:"page size"`
}
type CollectionGetListRes struct {
	Data       []*DocCollectionItem `json:"data"`
	Total      int                  `json:"total"`
	Page       int                  `json:"page"`
	PageSize   int                  `json:"page_size"`
	TotalPages int                  `json:"total_pages"`
}

// ----------------------------------------------------------------
//  Doc CRUD
// ----------------------------------------------------------------

type DocCreateReq struct {
	g.Meta        `path:"/docs" method:"post" tags:"Doc" summary:"Create doc"`
	CollectionId  int64       `v:"required|min:1"                            dc:"collection id"`
	ParentId      *int64      `v:"min:1"                                     dc:"parent doc id (for nested chapters)"`
	Title         string      `v:"required|length:1,200"                     dc:"doc title"`
	Slug          string      `v:"required|length:1,200|regex:^[a-z0-9-]+$"  dc:"URL slug"`
	Content       string      `                                              dc:"doc content"`
	Excerpt       string      `v:"max-length:500"                            dc:"short excerpt"`
	Status        int         `v:"required|in:1,2,3"                         dc:"1=draft 2=published 3=archived"`
	CommentStatus int         `v:"in:0,1"                                    dc:"0=closed 1=open"`
	Locale        string      `v:"length:2,10"                               dc:"e.g. zh-CN"`
	SortOrder     int         `                                              dc:"sort order within collection"`
	PublishedAt   *gtime.Time `                                              dc:"scheduled publish time"`
}
type DocCreateRes struct {
	Id int64 `json:"id"`
}

type DocUpdateReq struct {
	g.Meta        `path:"/docs/{id}" method:"put" tags:"Doc" summary:"Update doc"`
	Id            int64       `v:"required|min:1"                            dc:"doc id"`
	CollectionId  *int64      `v:"min:1"                                     dc:"collection id"`
	ParentId      *int64      `v:"min:0"                                     dc:"parent doc id (0 = clear)"`
	Title         *string     `v:"length:1,200"                              dc:"doc title"`
	Slug          *string     `v:"length:1,200|regex:^[a-z0-9-]+$"           dc:"URL slug"`
	Content       *string     `                                              dc:"doc content"`
	Excerpt       *string     `v:"max-length:500"                            dc:"short excerpt"`
	Status        *int        `v:"in:1,2,3"                                  dc:"doc status"`
	CommentStatus *int        `v:"in:0,1"                                    dc:"comment status"`
	Locale        *string     `v:"length:2,10"                               dc:"locale"`
	SortOrder     *int        `                                              dc:"sort order"`
	PublishedAt   *gtime.Time `                                              dc:"scheduled publish time"`
}
type DocUpdateRes struct{}

type DocDeleteReq struct {
	g.Meta `path:"/docs/{id}" method:"delete" tags:"Doc" summary:"Delete doc"`
	Id     int64 `v:"required|min:1" dc:"doc id"`
}
type DocDeleteRes struct{}

type DocGetOneReq struct {
	g.Meta `path:"/docs/{id}" method:"get" tags:"Doc" summary:"Get doc by id"`
	Id     int64 `v:"required|min:1" dc:"doc id"`
}
type DocGetOneRes struct {
	*DocDetailItem `dc:"doc detail"`
}

type DocGetBySlugReq struct {
	g.Meta `path:"/docs/slug/{slug}" method:"get" tags:"Doc" summary:"Get doc by slug"`
	Slug   string `v:"required|length:1,200" dc:"doc slug"`
}
type DocGetBySlugRes struct {
	*DocDetailItem `dc:"doc detail"`
}

type DocGetListReq struct {
	g.Meta       `path:"/docs" method:"get" tags:"Doc" summary:"List docs"`
	CollectionId *int64  `p:"collection_id" v:"min:1"             dc:"filter by collection"`
	ParentId     *int64  `p:"parent_id" v:"min:0"                 dc:"filter by parent doc (0 = top level)"`
	Status       *int    `p:"status" v:"in:1,2,3"                dc:"filter by status"`
	AuthorId     *int64  `p:"author_id" v:"min:1"                dc:"filter by author"`
	Locale       *string `p:"locale" v:"length:2,10"             dc:"filter by locale"`
	Keyword      *string `p:"keyword"                            dc:"search title"`
	Page         int     `p:"page" v:"min:1" d:"1"               dc:"page number"`
	PageSize     int     `p:"page_size" v:"between:1,100" d:"20" dc:"page size"`
}
type DocGetListRes struct {
	Data       []*DocItem `json:"data"`
	Total      int        `json:"total"`
	Page       int        `json:"page"`
	PageSize   int        `json:"page_size"`
	TotalPages int        `json:"total_pages"`
}

// ----------------------------------------------------------------
//  Doc SEO
// ----------------------------------------------------------------

type DocSeoUpdateReq struct {
	g.Meta         `path:"/docs/{id}/seo" method:"put" tags:"Doc" summary:"Update doc SEO"`
	Id             int64  `v:"required|min:1"  dc:"doc id"`
	MetaTitle      string `v:"max-length:200"  dc:"SEO title"`
	MetaDesc       string `v:"max-length:500"  dc:"SEO description"`
	OgTitle        string `v:"max-length:200"  dc:"Open Graph title"`
	OgImage        string `v:"max-length:500"  dc:"Open Graph image URL"`
	CanonicalUrl   string `v:"max-length:500"  dc:"canonical URL"`
	Robots         string `v:"max-length:100"  dc:"e.g. index,follow"`
	StructuredData string `                    dc:"JSON-LD structured data"`
}
type DocSeoUpdateRes struct{}

// ----------------------------------------------------------------
//  Doc Revisions
// ----------------------------------------------------------------

type DocRevisionListReq struct {
	g.Meta `path:"/docs/{id}/revisions" method:"get" tags:"Doc" summary:"List doc revisions"`
	Id     int64 `v:"required|min:1" dc:"doc id"`
	Page   int   `v:"min:1" d:"1"    dc:"page number"`
	Size   int   `v:"between:1,50" d:"10" dc:"page size"`
}
type DocRevisionListRes struct {
	List  []*DocRevisionItem `json:"list"`
	Total int                `json:"total"`
}

type DocRevisionRestoreReq struct {
	g.Meta     `path:"/docs/{id}/revisions/{revision_id}/restore" method:"post" tags:"Doc" summary:"Restore doc to revision"`
	Id         int64 `v:"required|min:1" dc:"doc id"`
	RevisionId int64 `v:"required|min:1" dc:"revision id"`
}
type DocRevisionRestoreRes struct{}

// ----------------------------------------------------------------
//  Doc view
// ----------------------------------------------------------------

type DocViewReq struct {
	g.Meta `path:"/docs/{id}/view" method:"post" tags:"Doc" summary:"Increment doc view count"`
	Id     int64 `v:"required|min:1" dc:"doc id"`
}
type DocViewRes struct{}
