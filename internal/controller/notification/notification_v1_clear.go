package notification

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NotificationClear(ctx context.Context, req *v1.NotificationClearReq) (res *v1.NotificationClearRes, err error) {
	return &v1.NotificationClearRes{}, service.Notification().Clear(ctx, req.UserId, req.Filter)
}
