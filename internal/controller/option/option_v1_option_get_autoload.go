package option

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) OptionGetAutoload(ctx context.Context, req *v1.OptionGetAutoloadReq) (res *v1.OptionGetAutoloadRes, err error) {
	options, err := service.Option().GetAutoload(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.OptionGetAutoloadRes{Options: options}, nil
}
