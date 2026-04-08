package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/notification/v1"
	"context"
)

type (
	INotification interface {
		List(ctx context.Context, req *v1.NotificationListReq) (*v1.NotificationListRes, error)
		UnreadCount(ctx context.Context, userId int64) (int, error)
		MarkRead(ctx context.Context, id int64) error
		MarkAllRead(ctx context.Context, userId int64) error
		Delete(ctx context.Context, id int64) error
		Clear(ctx context.Context, userId int64, filter string) error
		Create(ctx context.Context, notifType, subType string, actorId *int64, actorName, actorAvatar string, userId int64, objectType string, objectId *int64, objectTitle, objectLink, content string) error
	}
)

var localNotification INotification

func Notification() INotification {
	if localNotification == nil {
		panic("implement not found for interface INotification, forgot register?")
	}
	return localNotification
}

func RegisterNotification(i INotification) {
	localNotification = i
}
