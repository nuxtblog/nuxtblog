package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaGetList(ctx context.Context, req *v1.MediaGetListReq) (res *v1.MediaGetListRes, err error) {
	return service.Media().GetList(ctx, req)
}
