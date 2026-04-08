package user

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) UserGetOne(ctx context.Context, req *v1.UserGetOneReq) (res *v1.UserGetOneRes, err error) {
	item, err := service.User().GetOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.user_not_found"))
	}
	return &v1.UserGetOneRes{UserItem: item}, nil
}
