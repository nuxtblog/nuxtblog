package notification

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NotificationDelete(ctx context.Context, req *v1.NotificationDeleteReq) (res *v1.NotificationDeleteRes, err error) {
	return &v1.NotificationDeleteRes{}, service.Notification().Delete(ctx, req.Id)
}
