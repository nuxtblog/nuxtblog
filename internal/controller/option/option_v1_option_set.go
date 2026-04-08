package option

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) OptionSet(ctx context.Context, req *v1.OptionSetReq) (res *v1.OptionSetRes, err error) {
	return &v1.OptionSetRes{}, service.Option().Set(ctx, req)
}
