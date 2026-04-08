package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaUpdate(ctx context.Context, req *v1.MediaUpdateReq) (res *v1.MediaUpdateRes, err error) {
	return &v1.MediaUpdateRes{}, service.Media().Update(ctx, req)
}
