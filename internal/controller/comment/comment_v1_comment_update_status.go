package comment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CommentUpdateStatus(ctx context.Context, req *v1.CommentUpdateStatusReq) (res *v1.CommentUpdateStatusRes, err error) {
	return &v1.CommentUpdateStatusRes{}, service.Comment().UpdateStatus(ctx, req.Id, req.Status)
}
