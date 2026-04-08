package notification

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) NotificationUnreadCount(ctx context.Context, req *v1.NotificationUnreadCountReq) (res *v1.NotificationUnreadCountRes, err error) {
	count, err := service.Notification().UnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &v1.NotificationUnreadCountRes{Count: count}, nil
}
