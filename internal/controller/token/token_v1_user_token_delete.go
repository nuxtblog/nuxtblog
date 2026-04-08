package token

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/token/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserTokenDelete(ctx context.Context, req *v1.UserTokenDeleteReq) (res *v1.UserTokenDeleteRes, err error) {
	return service.Token().Delete(ctx, req)
}
