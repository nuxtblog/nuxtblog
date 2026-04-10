package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

type MomentVisibility int

const (
	MomentVisibilityPublic    MomentVisibility = 1
	MomentVisibilityPrivate   MomentVisibility = 2
	MomentVisibilityFollowers MomentVisibility = 3
)

type MomentMediaItem struct {
	Id       int64  `json:"id"`
	Url      string `json:"url"`
	MimeType string `json:"mime_type"`
	Width    int    `json:"width,omitempty"`
	Height   int    `json:"height,omitempty"`
}

type MomentAuthorItem struct {
	Id       int64  `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar,omitempty"`
}

type MomentStatsItem struct {
	ViewCount    int64 `json:"view_count"`
	LikeCount    int64 `json:"like_count"`
	CommentCount int64 `json:"comment_count"`
}

type MomentItem struct {
	Id         int64              `json:"id"`
	AuthorId   int64              `json:"author_id"`
	Content    string             `json:"content"`
	Visibility int                `json:"visibility"`
	CreatedAt  *gtime.Time        `json:"created_at"`
	UpdatedAt  *gtime.Time        `json:"updated_at"`
	Author     *MomentAuthorItem  `json:"author,omitempty"`
	Media      []*MomentMediaItem `json:"media,omitempty"`
	Stats      *MomentStatsItem   `json:"stats,omitempty"`
}

// ----------------------------------------------------------------
//  Create
// ----------------------------------------------------------------

type MomentCreateReq struct {
	g.Meta     `path:"/moments" method:"post" tags:"Moment" summary:"Create moment"`
	Content    string  `v:"required|length:1,5000"     dc:"moment content"`
	Visibility int     `v:"in:1,2,3" d:"1"             dc:"1=public 2=private 3=followers"`
	MediaIds   []int64 `p:"media_ids"                  dc:"ordered list of media ids to attach"`
}
type MomentCreateRes struct {
	Id int64 `json:"id"`
}

// ----------------------------------------------------------------
//  Update
// ----------------------------------------------------------------

type MomentUpdateReq struct {
	g.Meta     `path:"/moments/{id}" method:"put" tags:"Moment" summary:"Update moment"`
	Id         int64    `v:"required|min:1"            dc:"moment id"`
	Content    *string  `v:"length:1,5000"             dc:"moment content"`
	Visibility *int     `v:"in:1,2,3"                  dc:"visibility"`
	MediaIds   *[]int64 `p:"media_ids"                 dc:"replace media attachments (nil = no change)"`
}
type MomentUpdateRes struct{}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type MomentDeleteReq struct {
	g.Meta `path:"/moments/{id}" method:"delete" tags:"Moment" summary:"Delete moment"`
	Id     int64 `v:"required|min:1" dc:"moment id"`
}
type MomentDeleteRes struct{}

// ----------------------------------------------------------------
//  Get one
// ----------------------------------------------------------------

type MomentGetOneReq struct {
	g.Meta `path:"/moments/{id}" method:"get" tags:"Moment" summary:"Get moment by id"`
	Id     int64 `v:"required|min:1" dc:"moment id"`
}
type MomentGetOneRes struct {
	*MomentItem `dc:"moment detail"`
}

// ----------------------------------------------------------------
//  Get list
// ----------------------------------------------------------------

type MomentGetListReq struct {
	g.Meta     `path:"/moments" method:"get" tags:"Moment" summary:"List moments"`
	AuthorId   *int64  `p:"author_id" v:"min:1"              dc:"filter by author"`
	Visibility *int    `p:"visibility" v:"in:1,2,3"          dc:"filter by visibility"`
	Keyword    *string `p:"keyword"                          dc:"search content"`
	Page       int     `p:"page" v:"min:1" d:"1"             dc:"page number"`
	PageSize   int     `p:"page_size" v:"between:1,100" d:"20" dc:"page size"`
}
type MomentGetListRes struct {
	Data       []*MomentItem `json:"data"`
	Total      int           `json:"total"`
	Page       int           `json:"page"`
	PageSize   int           `json:"page_size"`
	TotalPages int           `json:"total_pages"`
}

// ----------------------------------------------------------------
//  View
// ----------------------------------------------------------------

type MomentViewReq struct {
	g.Meta `path:"/moments/{id}/view" method:"post" tags:"Moment" summary:"Increment moment view count"`
	Id     int64 `v:"required|min:1" dc:"moment id"`
}
type MomentViewRes struct{}
