package user

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) UserGetList(ctx context.Context, req *v1.UserGetListReq) (res *v1.UserGetListRes, err error) {
	return service.User().GetList(ctx, req)
}
