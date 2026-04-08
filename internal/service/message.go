package service

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
)

type IMessage interface {
	ConversationList(ctx context.Context, req *v1.ConversationListReq) (*v1.ConversationListRes, error)
	Send(ctx context.Context, req *v1.MessageSendReq) (*v1.MessageSendRes, error)
	MessageList(ctx context.Context, req *v1.MessageListReq) (*v1.MessageListRes, error)
	UnreadCount(ctx context.Context) (int, error)
}

var _message IMessage

func Message() IMessage         { return _message }
func RegisterMessage(m IMessage) { _message = m }
