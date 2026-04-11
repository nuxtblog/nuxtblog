package v1

import "github.com/gogf/gf/v2/frame/g"

// ── Shared item ──────────────────────────────────────────────────────────────

type FriendlinkItem struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Url         string `json:"url"`
	Logo        string `json:"logo"`
	Description string `json:"description"`
	SortOrder   int    `json:"sort_order"`
	Status      int    `json:"status"`
	CreatedAt   string `json:"created_at"`
}

// ── Admin: create ────────────────────────────────────────────────────────────

type FriendlinkCreateReq struct {
	g.Meta      `path:"/admin/friendlinks" method:"post" tags:"Friendlink" summary:"Admin: create friendlink"`
	Name        string `v:"required|length:1,200"  dc:"site name"`
	Url         string `v:"required|url"           dc:"site url"`
	Logo        string `dc:"site logo url"`
	Description string `dc:"short description"`
	SortOrder   int    `d:"0"  dc:"sort order (lower first)"`
	Status      int    `d:"1"  v:"in:0,1" dc:"1=visible 0=hidden"`
}
type FriendlinkCreateRes struct {
	Id int64 `json:"id"`
}

// ── Admin: list ──────────────────────────────────────────────────────────────

type FriendlinkAdminListReq struct {
	g.Meta `path:"/admin/friendlinks" method:"get" tags:"Friendlink" summary:"Admin: list friendlinks"`
	Page   int `d:"1"  v:"min:1"`
	Size   int `d:"20" v:"min:1|max:100"`
}
type FriendlinkAdminListRes struct {
	List  []*FriendlinkItem `json:"list"`
	Total int               `json:"total"`
}

// ── Admin: update ────────────────────────────────────────────────────────────

type FriendlinkUpdateReq struct {
	g.Meta      `path:"/admin/friendlinks/{id}" method:"put" tags:"Friendlink" summary:"Admin: update friendlink"`
	Id          int64  `v:"required|min:1"         dc:"friendlink id"`
	Name        string `v:"required|length:1,200"  dc:"site name"`
	Url         string `v:"required|url"           dc:"site url"`
	Logo        string `dc:"site logo url"`
	Description string `dc:"short description"`
	SortOrder   int    `dc:"sort order"`
	Status      int    `v:"in:0,1" dc:"1=visible 0=hidden"`
}
type FriendlinkUpdateRes struct{}

// ── Admin: delete ────────────────────────────────────────────────────────────

type FriendlinkDeleteReq struct {
	g.Meta `path:"/admin/friendlinks/{id}" method:"delete" tags:"Friendlink" summary:"Admin: delete friendlink"`
	Id     int64 `v:"required|min:1" dc:"friendlink id"`
}
type FriendlinkDeleteRes struct{}

// ── Public: list ─────────────────────────────────────────────────────────────

type FriendlinkListReq struct {
	g.Meta `path:"/friendlinks" method:"get" tags:"Friendlink" summary:"Public: list visible friendlinks"`
}
type FriendlinkListRes struct {
	List []*FriendlinkItem `json:"list"`
}
