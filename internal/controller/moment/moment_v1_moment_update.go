package moment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MomentUpdate(ctx context.Context, req *v1.MomentUpdateReq) (res *v1.MomentUpdateRes, err error) {
	if err = service.Moment().Update(ctx, req); err != nil {
		return nil, err
	}
	return &v1.MomentUpdateRes{}, nil
}
