package auth

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AuthLogout(ctx context.Context, req *v1.AuthLogoutReq) (res *v1.AuthLogoutRes, err error) {
	return &v1.AuthLogoutRes{}, service.Auth().Logout(ctx, req.RefreshToken)
}
