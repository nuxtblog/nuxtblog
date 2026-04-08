// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// Messages is the golang structure for table messages.
type Messages struct {
	Id             int64       `json:"id"             orm:"id"              description:""`
	ConversationId int64       `json:"conversationId" orm:"conversation_id" description:""`
	SenderId       int64       `json:"senderId"       orm:"sender_id"       description:""`
	Content        string      `json:"content"        orm:"content"         description:""`
	IsRead         int         `json:"isRead"         orm:"is_read"         description:""`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
}
