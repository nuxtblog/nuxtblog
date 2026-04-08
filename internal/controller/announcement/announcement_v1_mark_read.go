package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *PublicControllerV1) AnnouncementMarkRead(ctx context.Context, req *v1.AnnouncementMarkReadReq) (*v1.AnnouncementMarkReadRes, error) {
	userId, _ := middleware.GetCurrentUserID(ctx)
	return &v1.AnnouncementMarkReadRes{}, service.Announcement().MarkRead(ctx, userId)
}
