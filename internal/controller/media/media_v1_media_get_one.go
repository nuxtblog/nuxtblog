package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) MediaGetOne(ctx context.Context, req *v1.MediaGetOneReq) (res *v1.MediaGetOneRes, err error) {
	item, err := service.Media().GetOne(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, g.I18n().T(ctx, "error.media_not_found"))
	}
	return &v1.MediaGetOneRes{MediaItem: item}, nil
}
