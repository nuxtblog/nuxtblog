package comment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CommentDelete(ctx context.Context, req *v1.CommentDeleteReq) (res *v1.CommentDeleteRes, err error) {
	return &v1.CommentDeleteRes{}, service.Comment().Delete(ctx, req.Id)
}
