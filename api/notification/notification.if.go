// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package notification

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/notification/v1"
)

type INotificationV1 interface {
	NotificationList(ctx context.Context, req *v1.NotificationListReq) (res *v1.NotificationListRes, err error)
	NotificationUnreadCount(ctx context.Context, req *v1.NotificationUnreadCountReq) (res *v1.NotificationUnreadCountRes, err error)
	NotificationRead(ctx context.Context, req *v1.NotificationReadReq) (res *v1.NotificationReadRes, err error)
	NotificationReadAll(ctx context.Context, req *v1.NotificationReadAllReq) (res *v1.NotificationReadAllRes, err error)
	NotificationDelete(ctx context.Context, req *v1.NotificationDeleteReq) (res *v1.NotificationDeleteRes, err error)
	NotificationClear(ctx context.Context, req *v1.NotificationClearReq) (res *v1.NotificationClearRes, err error)
}
