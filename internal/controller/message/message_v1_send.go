package message

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MessageSend(ctx context.Context, req *v1.MessageSendReq) (res *v1.MessageSendRes, err error) {
	return service.Message().Send(ctx, req)
}
