package moment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MomentGetList(ctx context.Context, req *v1.MomentGetListReq) (res *v1.MomentGetListRes, err error) {
	return service.Moment().GetList(ctx, req)
}
