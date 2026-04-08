package comment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CommentUpdateContent(ctx context.Context, req *v1.CommentUpdateContentReq) (res *v1.CommentUpdateContentRes, err error) {
	return &v1.CommentUpdateContentRes{}, service.Comment().UpdateContent(ctx, req.Id, req.Content)
}
