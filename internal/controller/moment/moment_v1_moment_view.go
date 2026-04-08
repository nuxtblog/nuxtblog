package moment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MomentView(ctx context.Context, req *v1.MomentViewReq) (res *v1.MomentViewRes, err error) {
	if err = service.Moment().IncrementView(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.MomentViewRes{}, nil
}
