// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import "github.com/gogf/gf/v2/os/gtime"

// Conversations is the golang structure for table conversations.
type Conversations struct {
	Id        int64       `json:"id"        orm:"id"          description:""`
	UserA     int64       `json:"userA"     orm:"user_a"      description:""`
	UserB     int64       `json:"userB"     orm:"user_b"      description:""`
	LastMsg   string      `json:"lastMsg"   orm:"last_msg"    description:""`
	LastMsgAt *gtime.Time `json:"lastMsgAt" orm:"last_msg_at" description:""`
	UnreadA   int         `json:"unreadA"   orm:"unread_a"    description:""`
	UnreadB   int         `json:"unreadB"   orm:"unread_b"    description:""`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"  description:""`
}
