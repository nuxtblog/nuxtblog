package entity

import "github.com/gogf/gf/v2/os/gtime"

type Docs struct {
	Id            int         `json:"id"            orm:"id"`
	CollectionId  int         `json:"collectionId"  orm:"collection_id"`
	ParentId      *int        `json:"parentId"      orm:"parent_id"`
	SortOrder     int         `json:"sortOrder"     orm:"sort_order"`
	Status        int         `json:"status"        orm:"status"`
	Title         string      `json:"title"         orm:"title"`
	Slug          string      `json:"slug"          orm:"slug"`
	Content       string      `json:"content"       orm:"content"`
	Excerpt       string      `json:"excerpt"       orm:"excerpt"`
	AuthorId      int         `json:"authorId"      orm:"author_id"`
	CommentStatus int         `json:"commentStatus" orm:"comment_status"`
	Locale        string      `json:"locale"        orm:"locale"`
	PublishedAt   *gtime.Time `json:"publishedAt"   orm:"published_at"`
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"`
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"`
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"`
}
