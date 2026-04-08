package checkin

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/checkin/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) DoCheckin(ctx context.Context, req *v1.DoCheckinReq) (res *v1.DoCheckinRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	already, streak, err := service.Checkin().DoCheckin(ctx, uid)
	if err != nil {
		return nil, err
	}
	return &v1.DoCheckinRes{
		AlreadyCheckedIn: already,
		Streak:           streak,
	}, nil
}
