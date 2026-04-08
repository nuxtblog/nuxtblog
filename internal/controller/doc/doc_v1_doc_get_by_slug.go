package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) DocGetBySlug(ctx context.Context, req *v1.DocGetBySlugReq) (res *v1.DocGetBySlugRes, err error) {
	item, err := service.Doc().DocGetBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "doc not found")
	}
	return &v1.DocGetBySlugRes{DocDetailItem: item}, nil
}
