package message

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/message/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ConversationList(ctx context.Context, req *v1.ConversationListReq) (res *v1.ConversationListRes, err error) {
	return service.Message().ConversationList(ctx, req)
}
