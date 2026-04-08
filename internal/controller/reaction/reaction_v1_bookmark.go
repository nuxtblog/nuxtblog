package reaction

import (
	"context"
	"errors"

	v1 "github.com/nuxtblog/nuxtblog/api/reaction/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) BookmarkPost(ctx context.Context, req *v1.BookmarkPostReq) (res *v1.BookmarkPostRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return &v1.BookmarkPostRes{}, service.Reaction().Bookmark(ctx, req.Id, uid)
}

func (c *ControllerV1) UnbookmarkPost(ctx context.Context, req *v1.UnbookmarkPostReq) (res *v1.UnbookmarkPostRes, err error) {
	uid, ok := middleware.GetCurrentUserID(ctx)
	if !ok {
		return nil, errors.New(g.I18n().T(ctx, "error.unauthorized"))
	}
	return &v1.UnbookmarkPostRes{}, service.Reaction().Unbookmark(ctx, req.Id, uid)
}
