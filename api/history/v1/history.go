package v1

import "github.com/gogf/gf/v2/frame/g"

type HistoryListReq struct {
	g.Meta `path:"/history" method:"get" tags:"History" summary:"Get browse history" auth:"true"`
	Page   int `d:"1"  v:"min:1"`
	Size   int `d:"20" v:"min:1|max:50"`
}

type HistoryItem struct {
	PostId    int64  `json:"post_id"`
	PostTitle string `json:"post_title"`
	PostSlug  string `json:"post_slug"`
	PostCover string `json:"post_cover"`
	ViewedAt  string `json:"viewed_at"`
}

type HistoryListRes struct {
	Items []HistoryItem `json:"items"`
	Total int           `json:"total"`
}

type HistoryClearReq struct {
	g.Meta `path:"/history" method:"delete" tags:"History" summary:"Clear browse history" auth:"true"`
}

type HistoryClearRes struct{}
