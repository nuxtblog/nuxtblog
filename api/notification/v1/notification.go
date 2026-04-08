package v1

import "github.com/gogf/gf/v2/frame/g"

// NotificationItem is the response item for notifications
type NotificationItem struct {
	Id           int64  `json:"id"`
	Type         string `json:"type"`                   // follow|like|comment|reply|mention|system
	SubType      string `json:"sub_type,omitempty"`     // system sub type
	UserName     string `json:"user_name,omitempty"`    // actor name
	Avatar       string `json:"avatar,omitempty"`       // actor avatar
	Action       string `json:"action,omitempty"`       // "评论了你的文章"
	Title        string `json:"title,omitempty"`        // system notification title
	Content      string `json:"content"`
	RelatedTitle string `json:"related_title,omitempty"`
	RelatedLink  string `json:"related_link,omitempty"`
	Read         bool   `json:"read"`
	CreatedAt    string `json:"created_at"`
}

// NotificationListReq
type NotificationListReq struct {
	g.Meta `path:"/notifications" method:"get" tags:"Notification" summary:"Get notification list"`
	UserId int64  `p:"user_id" v:"required|min:1" dc:"recipient user id"`
	Filter string `p:"filter" d:"all"             dc:"all|unread|interaction|system"`
	Page   int    `p:"page"   v:"min:1"   d:"1"   dc:"page"`
	Size   int    `p:"size"   v:"between:1,100" d:"20" dc:"page size"`
}

type NotificationListRes struct {
	List       []*NotificationItem `json:"list"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	Size       int                 `json:"size"`
	TotalPages int                 `json:"total_pages"`
	Unread     int                 `json:"unread"`
}

// NotificationUnreadCountReq
type NotificationUnreadCountReq struct {
	g.Meta `path:"/notifications/unread-count" method:"get" tags:"Notification" summary:"Get unread count"`
	UserId int64 `p:"user_id" v:"required|min:1" dc:"recipient user id"`
}

type NotificationUnreadCountRes struct {
	Count int `json:"count"`
}

// NotificationReadReq
type NotificationReadReq struct {
	g.Meta `path:"/notifications/{id}/read" method:"put" tags:"Notification" summary:"Mark notification as read"`
	Id     int64 `v:"required|min:1" dc:"notification id"`
}
type NotificationReadRes struct{}

// NotificationReadAllReq
type NotificationReadAllReq struct {
	g.Meta `path:"/notifications/read-all" method:"put" tags:"Notification" summary:"Mark all notifications as read"`
	UserId int64 `p:"user_id" v:"required|min:1" dc:"user id"`
}
type NotificationReadAllRes struct{}

// NotificationDeleteReq
type NotificationDeleteReq struct {
	g.Meta `path:"/notifications/{id}" method:"delete" tags:"Notification" summary:"Delete a notification"`
	Id     int64 `v:"required|min:1" dc:"notification id"`
}
type NotificationDeleteRes struct{}

// NotificationClearReq
type NotificationClearReq struct {
	g.Meta `path:"/notifications/clear" method:"delete" tags:"Notification" summary:"Clear notifications by filter"`
	UserId int64  `p:"user_id" v:"required|min:1" dc:"user id"`
	Filter string `p:"filter"  d:"all"            dc:"all|unread|interaction|system"`
}
type NotificationClearRes struct{}
