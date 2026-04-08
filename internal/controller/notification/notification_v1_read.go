package notification

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NotificationRead(ctx context.Context, req *v1.NotificationReadReq) (res *v1.NotificationReadRes, err error) {
	return &v1.NotificationReadRes{}, service.Notification().MarkRead(ctx, req.Id)
}
