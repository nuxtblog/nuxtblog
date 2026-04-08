package auth

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AuthRefresh(ctx context.Context, req *v1.AuthRefreshReq) (res *v1.AuthRefreshRes, err error) {
	return service.Auth().Refresh(ctx, req.RefreshToken)
}
