package auth

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AuthLogin(ctx context.Context, req *v1.AuthLoginReq) (res *v1.AuthLoginRes, err error) {
	return service.Auth().Login(ctx, req)
}
