package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *PublicControllerV1) AnnouncementList(ctx context.Context, req *v1.AnnouncementListReq) (*v1.AnnouncementListRes, error) {
	userId, _ := middleware.GetCurrentUserID(ctx)
	list, total, unread, err := service.Announcement().ListForUser(ctx, userId, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.AnnouncementListRes{
		List:        list,
		Total:       total,
		UnreadCount: unread,
	}, nil
}
