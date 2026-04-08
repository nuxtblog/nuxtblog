package user

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserChangePassword(ctx context.Context, req *v1.UserChangePasswordReq) (res *v1.UserChangePasswordRes, err error) {
	return &v1.UserChangePasswordRes{}, service.User().ChangePassword(ctx, req)
}
