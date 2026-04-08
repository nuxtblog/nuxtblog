package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
)

type IAnnouncement interface {
	Create(ctx context.Context, title, content, atype string, createdBy int64) (int64, error)
	ListAdmin(ctx context.Context, page, size int) ([]*v1.AnnouncementItem, int, error)
	Update(ctx context.Context, id int64, title, content, atype string) error
	Delete(ctx context.Context, id int64) error
	ListForUser(ctx context.Context, userId int64, page, size int) ([]*v1.AnnouncementItem, int, int, error)
	MarkRead(ctx context.Context, userId int64) error
	UnreadCount(ctx context.Context, userId int64) (int, error)
}

var localAnnouncement IAnnouncement

func Announcement() IAnnouncement {
	if localAnnouncement == nil {
		panic("implement not found for interface IAnnouncement, forgot register?")
	}
	return localAnnouncement
}

func RegisterAnnouncement(i IAnnouncement) {
	localAnnouncement = i
}
