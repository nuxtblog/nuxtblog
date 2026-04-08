package checkin

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/checkin/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetCheckinStatus(ctx context.Context, req *v1.GetCheckinStatusReq) (res *v1.GetCheckinStatusRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	today, streak, err := service.Checkin().GetStatus(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &v1.GetCheckinStatusRes{
		CheckedInToday: today,
		Streak:         streak,
	}, nil
}
