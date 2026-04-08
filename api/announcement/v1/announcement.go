package v1

import "github.com/gogf/gf/v2/frame/g"

// ── Shared item ──────────────────────────────────────────────────────────────

type AnnouncementItem struct {
	Id        int64  `json:"id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	Type      string `json:"type"`
	CreatedAt string `json:"created_at"`
	Unread    bool   `json:"unread"`
}

// ── Admin: create ─────────────────────────────────────────────────────────────

type AnnouncementCreateReq struct {
	g.Meta  `path:"/admin/announcements" method:"post" tags:"Announcement" summary:"Admin: create announcement"`
	Title   string `v:"required|length:1,200"             dc:"announcement title"`
	Content string `v:"required|length:1,5000"            dc:"announcement body"`
	Type    string `v:"in:info,warning,success,danger" d:"info" dc:"display type"`
}
type AnnouncementCreateRes struct {
	Id int64 `json:"id"`
}

// ── Admin: list ───────────────────────────────────────────────────────────────

type AnnouncementListAdminReq struct {
	g.Meta `path:"/admin/announcements" method:"get" tags:"Announcement" summary:"Admin: list announcements"`
	Page   int `d:"1"  v:"min:1"`
	Size   int `d:"20" v:"min:1|max:100"`
}
type AnnouncementListAdminRes struct {
	List  []*AnnouncementItem `json:"list"`
	Total int                 `json:"total"`
}

// ── Admin: update ─────────────────────────────────────────────────────────────

type AnnouncementUpdateReq struct {
	g.Meta  `path:"/admin/announcements/{id}" method:"put" tags:"Announcement" summary:"Admin: update announcement"`
	Id      int64  `v:"required|min:1" dc:"announcement id"`
	Title   string `v:"required|length:1,200"          dc:"announcement title"`
	Content string `v:"required|length:1,5000"         dc:"announcement body"`
	Type    string `v:"in:info,warning,success,danger" dc:"display type"`
}
type AnnouncementUpdateRes struct{}

// ── Admin: delete ─────────────────────────────────────────────────────────────

type AnnouncementDeleteReq struct {
	g.Meta `path:"/admin/announcements/{id}" method:"delete" tags:"Announcement" summary:"Admin: delete announcement"`
	Id     int64 `v:"required|min:1" dc:"announcement id"`
}
type AnnouncementDeleteRes struct{}

// ── Public (authed): list ─────────────────────────────────────────────────────

type AnnouncementListReq struct {
	g.Meta `path:"/announcements" method:"get" tags:"Announcement" summary:"List announcements"`
	Page   int `d:"1"  v:"min:1"`
	Size   int `d:"20" v:"min:1|max:100"`
}
type AnnouncementListRes struct {
	List        []*AnnouncementItem `json:"list"`
	Total       int                 `json:"total"`
	UnreadCount int                 `json:"unread_count"`
}

// ── Public (authed): mark all read ────────────────────────────────────────────

type AnnouncementMarkReadReq struct {
	g.Meta `path:"/announcements/read" method:"put" tags:"Announcement" summary:"Mark all announcements as read (cursor update)"`
}
type AnnouncementMarkReadRes struct{}
