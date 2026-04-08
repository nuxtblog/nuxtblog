package reaction

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/reaction/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) GetReaction(ctx context.Context, req *v1.GetReactionReq) (res *v1.GetReactionRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return service.Reaction().GetStatus(ctx, req.Id, uid)
}

func (c *ControllerV1) GetBookmarks(ctx context.Context, req *v1.GetBookmarksReq) (res *v1.GetBookmarksRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return service.Reaction().GetBookmarks(ctx, uid, req.Page, req.Size)
}
