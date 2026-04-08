package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AdminControllerV1) AnnouncementDelete(ctx context.Context, req *v1.AnnouncementDeleteReq) (*v1.AnnouncementDeleteRes, error) {
	return &v1.AnnouncementDeleteRes{}, service.Announcement().Delete(ctx, req.Id)
}
