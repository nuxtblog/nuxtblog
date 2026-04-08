package user

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserCreate(ctx context.Context, req *v1.UserCreateReq) (res *v1.UserCreateRes, err error) {
	id, err := service.User().Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.UserCreateRes{Id: id}, nil
}
