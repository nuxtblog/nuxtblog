package doc

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/doc/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) DocGetOne(ctx context.Context, req *v1.DocGetOneReq) (res *v1.DocGetOneRes, err error) {
	item, err := service.Doc().DocGetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "doc not found")
	}
	return &v1.DocGetOneRes{DocDetailItem: item}, nil
}
