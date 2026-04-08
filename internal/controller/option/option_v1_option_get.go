package option

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) OptionGet(ctx context.Context, req *v1.OptionGetReq) (res *v1.OptionGetRes, err error) {
	item, err := service.Option().Get(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	return &v1.OptionGetRes{OptionItem: item}, nil
}
