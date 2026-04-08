package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) PostGetBySlug(ctx context.Context, req *v1.PostGetBySlugReq) (res *v1.PostGetBySlugRes, err error) {
	item, err := service.Post().GetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.post_not_found"))
	}
	return &v1.PostGetBySlugRes{PostDetailEnrichedItem: item}, nil
}
