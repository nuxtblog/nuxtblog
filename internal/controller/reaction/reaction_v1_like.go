package reaction

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/reaction/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) LikePost(ctx context.Context, req *v1.LikePostReq) (res *v1.LikePostRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return &v1.LikePostRes{}, service.Reaction().Like(ctx, req.Id, uid)
}

func (c *ControllerV1) UnlikePost(ctx context.Context, req *v1.UnlikePostReq) (res *v1.UnlikePostRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return &v1.UnlikePostRes{}, service.Reaction().Unlike(ctx, req.Id, uid)
}
