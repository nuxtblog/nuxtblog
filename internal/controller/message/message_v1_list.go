package message

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MessageList(ctx context.Context, req *v1.MessageListReq) (res *v1.MessageListRes, err error) {
	return service.Message().MessageList(ctx, req)
}
