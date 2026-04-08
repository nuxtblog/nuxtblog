package moment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MomentCreate(ctx context.Context, req *v1.MomentCreateReq) (res *v1.MomentCreateRes, err error) {
	id, err := service.Moment().Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.MomentCreateRes{Id: id}, nil
}
