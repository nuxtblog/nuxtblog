package user

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserUpdate(ctx context.Context, req *v1.UserUpdateReq) (res *v1.UserUpdateRes, err error) {
	return &v1.UserUpdateRes{}, service.User().Update(ctx, req)
}
