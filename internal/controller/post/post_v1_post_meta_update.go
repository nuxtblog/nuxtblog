package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// PostMetaUpdate 专用接口：PUT /posts/{id}/metas
func (c *ControllerV1) PostMetaUpdate(ctx context.Context, req *v1.PostMetaUpdateReq) (res *v1.PostMetaUpdateRes, err error) {
	detail, err := service.Post().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.post_not_found"))
	}
	return &v1.PostMetaUpdateRes{}, service.Post().UpsertMetas(ctx, req.Id, req.Metas)
}

// PostMetaGet 专用接口：GET /posts/{id}/metas
func (c *ControllerV1) PostMetaGet(ctx context.Context, req *v1.PostMetaGetReq) (res *v1.PostMetaGetRes, err error) {
	detail, err := service.Post().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.post_not_found"))
	}
	metas, err := service.Post().GetMetas(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.PostMetaGetRes{Metas: metas}, nil
}
