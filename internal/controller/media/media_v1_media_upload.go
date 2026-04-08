package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaUpload(ctx context.Context, req *v1.MediaUploadReq) (res *v1.MediaUploadRes, err error) {
	item, err := service.Media().Upload(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.MediaUploadRes{MediaItem: item}, nil
}
