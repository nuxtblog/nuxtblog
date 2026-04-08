package notification

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NotificationList(ctx context.Context, req *v1.NotificationListReq) (res *v1.NotificationListRes, err error) {
	return service.Notification().List(ctx, req)
}
