package announcement

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/announcement/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *AdminControllerV1) AnnouncementCreate(ctx context.Context, req *v1.AnnouncementCreateReq) (*v1.AnnouncementCreateRes, error) {
	createdBy, _ := middleware.GetCurrentUserID(ctx)
	id, err := service.Announcement().Create(ctx, req.Title, req.Content, req.Type, createdBy)
	if err != nil {
		return nil, err
	}
	return &v1.AnnouncementCreateRes{Id: id}, nil
}
