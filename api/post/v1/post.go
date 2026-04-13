package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ----------------------------------------------------------------
//  Enums
// ----------------------------------------------------------------

type PostType int

const (
	PostTypePost   PostType = 1
	PostTypePage   PostType = 2
	PostTypeCustom PostType = 3
)

type PostStatus int

const (
	PostStatusDraft     PostStatus = 1
	PostStatusPublished PostStatus = 2
	PostStatusPrivate   PostStatus = 3
	PostStatusArchived  PostStatus = 4
	PostStatusTrashed   PostStatus = 5
)

type CommentStatus int

const (
	CommentStatusOpen   CommentStatus = 1
	CommentStatusClosed CommentStatus = 0
)

// ----------------------------------------------------------------
//  Shared output
// ----------------------------------------------------------------

type PostItem struct {
	Id            int64         `json:"id"`
	PostType      PostType      `json:"post_type"`
	Status        PostStatus    `json:"status"`
	Title         string        `json:"title"`
	Slug          string        `json:"slug"`
	Excerpt       string        `json:"excerpt"`
	AuthorId      int64         `json:"author_id"`
	FeaturedImgId *int64        `json:"featured_img_id"`
	CommentStatus CommentStatus `json:"comment_status"`
	Locale        string        `json:"locale"`
	PublishedAt   *gtime.Time   `json:"published_at"`
	CreatedAt     *gtime.Time   `json:"created_at"`
	UpdatedAt     *gtime.Time   `json:"updated_at"`
	// 按需 preload
	// Author *User          `json:"author,omitempty"`
	Stats *PostStatsItem `json:"stats,omitempty"`
}

type PostDetailItem struct {
	PostItem
	Content string       `json:"content"`
	Seo     *PostSeoItem `json:"seo,omitempty"`
}

type PostStatsItem struct {
	ViewCount    int64 `json:"view_count"`
	LikeCount    int64 `json:"like_count"`
	CommentCount int64 `json:"comment_count"`
	ShareCount   int64 `json:"share_count"`
}

type PostSeoItem struct {
	MetaTitle      string `json:"meta_title"`
	MetaDesc       string `json:"meta_desc"`
	OgTitle        string `json:"og_title"`
	OgImage        string `json:"og_image"`
	CanonicalUrl   string `json:"canonical_url"`
	Robots         string `json:"robots"`
	StructuredData string `json:"structured_data"`
}

// ----------------------------------------------------------------
//  Create
// ----------------------------------------------------------------

type PostCreateReq struct {
	g.Meta          `path:"/posts" method:"post" tags:"Post" summary:"Create post"`
	PostType        PostType          `v:"required|in:1,2,3"          dc:"1=post 2=page 3=custom"`
	Title           string            `v:"required|length:1,200"       dc:"post title"`
	Slug            string            `v:"required|length:1,200|regex:^[a-z0-9-]+$" dc:"URL slug, lowercase letters, numbers, hyphens"`
	Content         string            `                                dc:"post content"`
	Excerpt         string            `v:"max-length:500"              dc:"short excerpt"`
	Status          PostStatus        `v:"required|in:1,2,3,4"         dc:"1=draft 2=published 3=private 4=archived"`
	CommentStatus   CommentStatus     `v:"in:0,1"                      dc:"0=closed 1=open"`
	FeaturedImgId   *int64            `v:"min:1"                       dc:"featured image media id"`
	Password        string            `v:"max-length:64"               dc:"password for protected post"`
	Locale          string            `v:"length:2,10"                 dc:"e.g. zh-CN"`
	PublishedAt     *gtime.Time       `                                dc:"scheduled publish time"`
	AuthorId        *int64            `p:"author_id" v:"min:1"         dc:"author user id (admin only, defaults to current user)"`
	TermTaxonomyIds []int64           `p:"term_taxonomy_ids"           dc:"list of term_taxonomy ids to associate"`
	Metas           map[string]string `p:"metas"                       dc:"key-value metadata, e.g. {\"is_sticky\":\"1\"}"`
}
type PostCreateRes struct {
	Id int64 `json:"id"`
}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type PostDeleteReq struct {
	g.Meta `path:"/posts/{id}" method:"delete" tags:"Post" summary:"Delete post"`
	Id     int64 `v:"required|min:1" dc:"post id"`
}
type PostDeleteRes struct{}

// ----------------------------------------------------------------
//  Update
// ----------------------------------------------------------------

type PostUpdateReq struct {
	g.Meta          `path:"/posts/{id}" method:"put" tags:"Post" summary:"Update post"`
	Id              int64             `v:"required|min:1"              dc:"post id"`
	Title           *string           `v:"length:1,200"                dc:"post title"`
	Slug            *string           `v:"length:1,200|regex:^[a-z0-9-]+$" dc:"URL slug"`
	Content         *string           `                                dc:"post content"`
	Excerpt         *string           `v:"max-length:500"              dc:"short excerpt"`
	Status          *PostStatus       `v:"in:1,2,3,4"                  dc:"post status"`
	CommentStatus   *CommentStatus    `v:"in:0,1"                      dc:"comment status"`
	FeaturedImgId   *int64            `v:"min:1"                       dc:"featured image media id"`
	Password        *string           `v:"max-length:64"               dc:"password"`
	Locale          *string           `v:"length:2,10"                 dc:"locale"`
	PublishedAt     *gtime.Time       `                                dc:"scheduled publish time"`
	AuthorId        *int64            `p:"author_id" v:"min:1"         dc:"change author (admin only)"`
	TermTaxonomyIds *[]int64          `p:"term_taxonomy_ids"           dc:"replace all term_taxonomy associations (nil = no change)"`
	Metas           map[string]string `p:"metas"                       dc:"upsert metas (nil=no change, empty map=clear all)"`
}
type PostUpdateRes struct{}

// ----------------------------------------------------------------
//  Get one  (by id or slug)
// ----------------------------------------------------------------

type PostGetOneReq struct {
	g.Meta `path:"/posts/{id}" method:"get" tags:"Post" summary:"Get post by id"`
	Id     int64 `v:"required|min:1" dc:"post id"`
}
type PostGetOneRes struct {
	*PostDetailItem `dc:"post detail"`
}

type PostGetBySlugReq struct {
	g.Meta `path:"/posts/slug/{slug}" method:"get" tags:"Post" summary:"Get post by slug"`
	Slug   string `v:"required|length:1,200" dc:"post slug"`
}
type PostGetBySlugRes struct {
	*PostDetailEnrichedItem
}

// PostDetailEnrichedItem is the enriched detail response for user-facing endpoints
type PostDetailEnrichedItem struct {
	Id            int64             `json:"id"`
	PostType      string            `json:"post_type"`
	Status        string            `json:"status"`
	Title         string            `json:"title"`
	Slug          string            `json:"slug"`
	Content       string            `json:"content"`
	Excerpt       string            `json:"excerpt"`
	CommentStatus string            `json:"comment_status"`
	Locale        string            `json:"locale"`
	PublishedAt   string            `json:"published_at,omitempty"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	ViewCount     int64             `json:"view_count"`
	CommentCount  int64             `json:"comment_count"`
	HasPassword    bool              `json:"has_password"`
	Author         *PostAuthorItem   `json:"author,omitempty"`
	FeaturedImg    *PostMediaItem    `json:"featured_img,omitempty"`
	Metas          map[string]string `json:"metas,omitempty"`
	// Paywall fields
	IsPaid         bool   `json:"is_paid"`
	Price          int    `json:"price,omitempty"`           // cents, 0=free
	PriceType      string `json:"price_type,omitempty"`      // "one_time" or "membership_only"
	IsUnlocked     bool   `json:"is_unlocked"`
	FreePreviewPct int    `json:"free_preview_pct,omitempty"` // 0-100
}

// ----------------------------------------------------------------
//  Get list
// ----------------------------------------------------------------

type PostGetListReq struct {
	g.Meta             `path:"/posts" method:"get" tags:"Post" summary:"Get post list"`
	PostType           *string `p:"post_type"                              dc:"post/page/custom"`
	Status             *string `p:"status"                                 dc:"draft/published/private/archived"`
	AuthorId           *int64  `p:"author_id" v:"min:1"                    dc:"filter by author"`
	Locale             *string `p:"locale" v:"length:2,10"                 dc:"filter by locale"`
	Keyword            *string `p:"keyword,search"                         dc:"search title"`
	TermTaxonomyId     *int64  `p:"term_taxonomy_id,tag" v:"min:1"         dc:"filter by taxonomy id"`
	IncludeCategoryIds string  `p:"include_category_ids"                   dc:"comma-separated term_taxonomy_ids to include (OR logic)"`
	ExcludeCategoryIds string  `p:"exclude_category_ids"                   dc:"comma-separated term_taxonomy_ids to exclude"`
	MetaKey            *string `p:"meta_key"                                dc:"filter by post meta key (requires exact match)"`
	MetaValue          *string `p:"meta_value"                              dc:"filter by meta value (only used when meta_key is set)"`
	SortBy             *string `p:"sort_by"                                 dc:"sort field: view_count (default: created_at desc)"`
	PublishedAfter     *string `p:"published_after"                         dc:"only posts published on or after this date (YYYY-MM-DD)"`
	Page               int     `p:"page" v:"min:1" d:"1"                   dc:"page number"`
	PageSize           int     `p:"page_size,size" v:"between:1,100" d:"20" dc:"page size"`
}
type PostGetListRes struct {
	Data       []*PostListItem `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// PostListItem is the enriched response item for list endpoints
type PostListItem struct {
	Id            int64             `json:"id"`
	PostType      string            `json:"post_type"`
	Status        string            `json:"status"`
	Title         string            `json:"title"`
	Slug          string            `json:"slug"`
	Excerpt       string            `json:"excerpt"`
	CommentStatus string            `json:"comment_status"`
	Locale        string            `json:"locale"`
	PublishedAt   string            `json:"published_at,omitempty"`
	CreatedAt     string            `json:"created_at"`
	UpdatedAt     string            `json:"updated_at"`
	ViewCount     int64             `json:"view_count"`
	CommentCount  int64             `json:"comment_count"`
	Author        *PostAuthorItem   `json:"author,omitempty"`
	FeaturedImg   *PostMediaItem    `json:"featured_img,omitempty"`
	Metas         map[string]string `json:"metas,omitempty"`
}

type PostAuthorItem struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type PostMediaItem struct {
	Id       int64  `json:"id"`
	Url      string `json:"url"`
	Title    string `json:"title,omitempty"`
	MimeType string `json:"mime_type,omitempty"`
}

// ----------------------------------------------------------------
//  SEO
// ----------------------------------------------------------------

type PostSeoUpdateReq struct {
	g.Meta         `path:"/posts/{id}/seo" method:"put" tags:"Post" summary:"Update post SEO"`
	Id             int64  `v:"required|min:1"   dc:"post id"`
	MetaTitle      string `v:"max-length:200"   dc:"SEO title"`
	MetaDesc       string `v:"max-length:500"   dc:"SEO description"`
	OgTitle        string `v:"max-length:200"   dc:"Open Graph title"`
	OgImage        string `v:"max-length:500"   dc:"Open Graph image URL"`
	CanonicalUrl   string `v:"max-length:500"   dc:"canonical URL"`
	Robots         string `v:"max-length:100"   dc:"e.g. index,follow"`
	StructuredData string `                     dc:"JSON-LD structured data"`
}
type PostSeoUpdateRes struct{}

// ----------------------------------------------------------------
//  Revisions
// ----------------------------------------------------------------

type PostRevisionListReq struct {
	g.Meta `path:"/posts/{id}/revisions" method:"get" tags:"Post" summary:"Get post revisions"`
	Id     int64 `v:"required|min:1" dc:"post id"`
	Page   int   `v:"min:1"          dc:"page number" d:"1"`
	Size   int   `v:"between:1,50"   dc:"page size"   d:"10"`
}

type PostRevisionItem struct {
	Id        int64       `json:"id"`
	PostId    int64       `json:"post_id"`
	AuthorId  int64       `json:"author_id"`
	Title     string      `json:"title"`
	RevNote   string      `json:"rev_note"`
	CreatedAt *gtime.Time `json:"created_at"`
}

type PostRevisionListRes struct {
	List  []*PostRevisionItem `json:"list"`
	Total int                 `json:"total"`
}

type PostRevisionRestoreReq struct {
	g.Meta     `path:"/posts/{id}/revisions/{revision_id}/restore" method:"post" tags:"Post" summary:"Restore post to revision"`
	Id         int64 `v:"required|min:1" dc:"post id"`
	RevisionId int64 `v:"required|min:1" dc:"revision id"`
}
type PostRevisionRestoreRes struct{}

// ----------------------------------------------------------------
//  Stats
// ----------------------------------------------------------------

type GetStatsReq struct {
	g.Meta `path:"/posts/stats" method:"get" tags:"Post" summary:"Get post statistics"`
}

type GetStatsRes struct {
	TotalPosts     int64 `json:"total_posts"`
	PublishedPosts int64 `json:"published_posts"`
	DraftPosts     int64 `json:"draft_posts"`
	PrivatePosts   int64 `json:"private_posts"`
	ArchivedPosts  int64 `json:"archived_posts"`
	TotalViews     int64 `json:"total_views"`
	TotalLikes     int64 `json:"total_likes"`
	TotalComments  int64 `json:"total_comments"`
}

// ----------------------------------------------------------------
//  Trash / Restore
// ----------------------------------------------------------------

type PostTrashReq struct {
	g.Meta `path:"/posts/{id}/trash" method:"post" tags:"Post" summary:"Move post to trash"`
	Id     int64 `v:"required|min:1" dc:"post id"`
}
type PostTrashRes struct{}

type PostRestoreReq struct {
	g.Meta `path:"/posts/{id}/restore" method:"post" tags:"Post" summary:"Restore post from trash"`
	Id     int64 `v:"required|min:1" dc:"post id"`
}
type PostRestoreRes struct{}

// ----------------------------------------------------------------
//  Batch
// ----------------------------------------------------------------

type PostBatchReq struct {
	g.Meta `path:"/posts/batch" method:"post" tags:"Post" summary:"Batch update posts"`
	Ids    []int64 `v:"required|min-length:1" dc:"post ids"`
	Action string  `v:"required|in:publish,draft,trash,restore,delete" dc:"batch action"`
}
type PostBatchRes struct {
	Affected int `json:"affected"`
}

// PostBatchUpdateReq applies field-level updates to multiple posts at once.
// Only non-nil fields are written. FeaturedImgId=0 clears the cover image.
// TermTaxonomyIds non-nil replaces all associations (empty slice = clear all).
type PostBatchUpdateReq struct {
	g.Meta          `path:"/posts/batch" method:"patch" tags:"Post" summary:"Batch update post fields"`
	Ids             []int64     `v:"required|min-length:1" dc:"post ids"`
	FeaturedImgId   *int64      `                           dc:"set featured image id; 0 = clear cover"`
	Status          *PostStatus `v:"in:1,2,3,4"            dc:"set status: 1=draft 2=published 3=private 4=archived"`
	TermTaxonomyIds *[]int64    `                           dc:"replace all taxonomy associations; empty slice = clear all"`
	AuthorId        *int64      `v:"min:1"                  dc:"change author (admin only)"`
}
type PostBatchUpdateRes struct {
	Affected int `json:"affected"`
}

// ----------------------------------------------------------------
//  View (increment view count)
// ----------------------------------------------------------------

type PostViewReq struct {
	g.Meta `path:"/posts/{id}/view" method:"post" tags:"Post" summary:"Increment post view count"`
	Id     int64 `v:"required|min:1" dc:"post id"`
}
type PostViewRes struct{}

// ----------------------------------------------------------------
//  Verify password
// ----------------------------------------------------------------

type PostVerifyPasswordReq struct {
	g.Meta   `path:"/posts/{id}/verify-password" method:"post" tags:"Post" summary:"Verify post password"`
	Id       int64  `v:"required|min:1" dc:"post id"`
	Password string `v:"required" dc:"post password to verify"`
}
type PostVerifyPasswordRes struct {
	Valid bool `json:"valid"`
}

// ----------------------------------------------------------------
//  Metas
// ----------------------------------------------------------------

// PostMetaUpdateReq upserts or deletes metas for a post.
// Keys present in the map are upserted; keys with empty string value are deleted.
type PostMetaUpdateReq struct {
	g.Meta `path:"/posts/{id}/metas" method:"put" tags:"Post" summary:"Upsert post metas"`
	Id     int64             `v:"required|min:1" dc:"post id"`
	Metas  map[string]string `v:"required"       dc:"map of meta_key→meta_value; empty value = delete that key"`
}
type PostMetaUpdateRes struct{}

type PostMetaGetReq struct {
	g.Meta `path:"/posts/{id}/metas" method:"get" tags:"Post" summary:"Get post metas"`
	Id     int64 `v:"required|min:1" dc:"post id"`
}
type PostMetaGetRes struct {
	Metas map[string]string `json:"metas"`
}
