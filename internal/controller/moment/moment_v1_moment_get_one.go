package moment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) MomentGetOne(ctx context.Context, req *v1.MomentGetOneReq) (res *v1.MomentGetOneRes, err error) {
	item, err := service.Moment().GetById(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "moment not found")
	}
	return &v1.MomentGetOneRes{MomentItem: item}, nil
}
