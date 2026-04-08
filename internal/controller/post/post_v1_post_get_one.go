package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) PostGetOne(ctx context.Context, req *v1.PostGetOneReq) (res *v1.PostGetOneRes, err error) {
	detail, err := service.Post().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.post_not_found"))
	}
	return &v1.PostGetOneRes{PostDetailItem: detail}, nil
}
