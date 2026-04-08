package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaCategoryList(ctx context.Context, req *v1.MediaCategoryListReq) (*v1.MediaCategoryListRes, error) {
	items, err := service.Media().GetCategories(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.MediaCategoryListRes{List: items}, nil
}
