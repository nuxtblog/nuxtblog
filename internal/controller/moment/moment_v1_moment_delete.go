package moment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MomentDelete(ctx context.Context, req *v1.MomentDeleteReq) (res *v1.MomentDeleteRes, err error) {
	if err = service.Moment().Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.MomentDeleteRes{}, nil
}
