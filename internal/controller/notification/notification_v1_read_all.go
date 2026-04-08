package notification

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NotificationReadAll(ctx context.Context, req *v1.NotificationReadAllReq) (res *v1.NotificationReadAllRes, err error) {
	return &v1.NotificationReadAllRes{}, service.Notification().MarkAllRead(ctx, req.UserId)
}
