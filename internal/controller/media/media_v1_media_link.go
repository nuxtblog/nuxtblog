package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaLink(ctx context.Context, req *v1.MediaLinkReq) (res *v1.MediaLinkRes, err error) {
	item, err := service.Media().Link(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.MediaLinkRes{MediaItem: item}, nil
}
