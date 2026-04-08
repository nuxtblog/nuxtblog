package message

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
)

type IMessage interface {
	ConversationList(ctx context.Context, req *v1.ConversationListReq) (res *v1.ConversationListRes, err error)
	MessageSend(ctx context.Context, req *v1.MessageSendReq) (res *v1.MessageSendRes, err error)
	MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error)
	MessageUnread(ctx context.Context, req *v1.MessageUnreadReq) (res *v1.MessageUnreadRes, err error)
}
