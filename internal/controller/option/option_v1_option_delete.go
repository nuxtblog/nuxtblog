package option

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) OptionDelete(ctx context.Context, req *v1.OptionDeleteReq) (res *v1.OptionDeleteRes, err error) {
	return &v1.OptionDeleteRes{}, service.Option().Delete(ctx, req.Key)
}
