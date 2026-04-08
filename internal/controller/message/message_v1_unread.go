package message

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MessageUnread(ctx context.Context, req *v1.MessageUnreadReq) (res *v1.MessageUnreadRes, err error) {
	count, err := service.Message().UnreadCount(ctx)
	return &v1.MessageUnreadRes{Count: count}, err
}
