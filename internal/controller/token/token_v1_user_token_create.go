package token

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/token/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserTokenCreate(ctx context.Context, req *v1.UserTokenCreateReq) (res *v1.UserTokenCreateRes, err error) {
	return service.Token().Create(ctx, req)
}
