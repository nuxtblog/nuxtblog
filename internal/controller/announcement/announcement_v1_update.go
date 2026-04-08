package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AdminControllerV1) AnnouncementUpdate(ctx context.Context, req *v1.AnnouncementUpdateReq) (*v1.AnnouncementUpdateRes, error) {
	err := service.Announcement().Update(ctx, req.Id, req.Title, req.Content, req.Type)
	if err != nil {
		return nil, err
	}
	return &v1.AnnouncementUpdateRes{}, nil
}
