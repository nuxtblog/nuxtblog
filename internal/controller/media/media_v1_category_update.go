package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaCategoryUpdate(ctx context.Context, req *v1.MediaCategoryUpdateReq) (*v1.MediaCategoryUpdateRes, error) {
	if err := service.Media().UpdateCategory(ctx, req); err != nil {
		return nil, err
	}
	return &v1.MediaCategoryUpdateRes{}, nil
}
