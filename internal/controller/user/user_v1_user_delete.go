package user

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/user/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) UserDelete(ctx context.Context, req *v1.UserDeleteReq) (res *v1.UserDeleteRes, err error) {
	r := ghttp.RequestFromCtx(ctx)
	if r != nil {
		if callerID := r.GetCtxVar("user_id").Int64(); callerID == req.Id {
			return nil, errors.New(g.I18n().T(ctx, "user.cannot_delete_self"))
		}
	}
	return &v1.UserDeleteRes{}, service.User().Delete(ctx, req.Id)
}
