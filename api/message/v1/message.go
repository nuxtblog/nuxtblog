package v1

import "github.com/gogf/gf/v2/frame/g"

type ConversationItem struct {
	Id          int64  `json:"id"`
	OtherUserId int64  `json:"other_user_id"`
	OtherName   string `json:"other_name"`
	OtherAvatar string `json:"other_avatar"`
	LastMsg     string `json:"last_msg"`
	LastMsgAt   string `json:"last_msg_at"`
	UnreadCount int    `json:"unread_count"`
}

type MessageItem struct {
	Id        int64  `json:"id"`
	SenderId  int64  `json:"sender_id"`
	Content   string `json:"content"`
	IsRead    bool   `json:"is_read"`
	CreatedAt string `json:"created_at"`
}

// List conversations
type ConversationListReq struct {
	g.Meta `path:"/messages" method:"get" tags:"Message" summary:"List my conversations" auth:"true"`
	Page   int `d:"1" v:"min:1"`
	Size   int `d:"20" v:"min:1|max:50"`
}
type ConversationListRes struct {
	Items       []ConversationItem `json:"items"`
	Total       int                `json:"total"`
	TotalUnread int                `json:"total_unread"`
}

// Send message
type MessageSendReq struct {
	g.Meta   `path:"/messages/{to_user_id}" method:"post" tags:"Message" summary:"Send a message" auth:"true"`
	ToUserId int64  `v:"required|min:1" dc:"recipient user id"`
	Content  string `v:"required|length:1,2000" dc:"message content"`
}
type MessageSendRes struct {
	Id             int64 `json:"id"`
	ConversationId int64 `json:"conversation_id"`
}

// Get messages in a conversation
type MessageListReq struct {
	g.Meta   `path:"/messages/{to_user_id}/history" method:"get" tags:"Message" summary:"Get conversation history" auth:"true"`
	ToUserId int64 `v:"required|min:1"`
	BeforeId int64 `d:"0" dc:"load messages before this id (for pagination)"`
	Size     int   `d:"30" v:"min:1|max:100"`
}
type MessageListRes struct {
	Items          []MessageItem `json:"items"`
	ConversationId int64         `json:"conversation_id"`
	HasMore        bool          `json:"has_more"`
}

// Total unread count (for header badge)
type MessageUnreadReq struct {
	g.Meta `path:"/messages/unread" method:"get" tags:"Message" summary:"Get total unread count" auth:"true"`
}
type MessageUnreadRes struct {
	Count int `json:"count"`
}
