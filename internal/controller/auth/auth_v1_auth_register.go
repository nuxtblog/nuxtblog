package auth

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AuthRegister(ctx context.Context, req *v1.AuthRegisterReq) (res *v1.AuthRegisterRes, err error) {
	return service.Auth().Register(ctx, req)
}
