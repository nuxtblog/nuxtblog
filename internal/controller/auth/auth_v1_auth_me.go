package auth

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/auth/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AuthMe(ctx context.Context, req *v1.AuthMeReq) (res *v1.AuthMeRes, err error) {
	user, err := service.Auth().Me(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.AuthMeRes{User: *user}, nil
}
