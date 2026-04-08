package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaLocalize(ctx context.Context, req *v1.MediaLocalizeReq) (res *v1.MediaLocalizeRes, err error) {
	item, err := service.Media().Localize(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.MediaLocalizeRes{MediaItem: item}, nil
}
