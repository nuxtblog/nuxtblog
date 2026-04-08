package entity

import "github.com/gogf/gf/v2/os/gtime"

type DocRevisions struct {
	Id        int         `json:"id"        orm:"id"`
	DocId     int         `json:"docId"     orm:"doc_id"`
	AuthorId  int         `json:"authorId"  orm:"author_id"`
	Title     string      `json:"title"     orm:"title"`
	Content   string      `json:"content"   orm:"content"`
	RevNote   string      `json:"revNote"   orm:"rev_note"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"`
}
