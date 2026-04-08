package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaDelete(ctx context.Context, req *v1.MediaDeleteReq) (res *v1.MediaDeleteRes, err error) {
	return &v1.MediaDeleteRes{}, service.Media().Delete(ctx, req.Id)
}
