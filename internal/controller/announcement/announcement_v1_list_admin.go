package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AdminControllerV1) AnnouncementListAdmin(ctx context.Context, req *v1.AnnouncementListAdminReq) (*v1.AnnouncementListAdminRes, error) {
	list, total, err := service.Announcement().ListAdmin(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	return &v1.AnnouncementListAdminRes{List: list, Total: total}, nil
}
