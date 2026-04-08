package token

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/token/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserTokenList(ctx context.Context, req *v1.UserTokenListReq) (res *v1.UserTokenListRes, err error) {
	return service.Token().List(ctx, req)
}
