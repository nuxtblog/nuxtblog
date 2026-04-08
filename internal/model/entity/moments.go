package entity

import "github.com/gogf/gf/v2/os/gtime"

type Moments struct {
	Id         int         `json:"id"         orm:"id"`
	AuthorId   int         `json:"authorId"   orm:"author_id"`
	Content    string      `json:"content"    orm:"content"`
	Visibility int         `json:"visibility" orm:"visibility"`
	CreatedAt  *gtime.Time `json:"createdAt"  orm:"created_at"`
	UpdatedAt  *gtime.Time `json:"updatedAt"  orm:"updated_at"`
	DeletedAt  *gtime.Time `json:"deletedAt"  orm:"deleted_at"`
}
