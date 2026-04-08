package entity

import "github.com/gogf/gf/v2/os/gtime"

type DocCollections struct {
	Id          int         `json:"id"          orm:"id"`
	Slug        string      `json:"slug"        orm:"slug"`
	Title       string      `json:"title"       orm:"title"`
	Description string      `json:"description" orm:"description"`
	CoverImgId  *int        `json:"coverImgId"  orm:"cover_img_id"`
	AuthorId    int         `json:"authorId"    orm:"author_id"`
	Status      int         `json:"status"      orm:"status"`
	Locale      string      `json:"locale"      orm:"locale"`
	SortOrder   int         `json:"sortOrder"   orm:"sort_order"`
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"`
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"`
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"`
}
