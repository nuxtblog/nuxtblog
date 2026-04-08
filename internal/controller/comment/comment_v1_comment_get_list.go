package comment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CommentGetList(ctx context.Context, req *v1.CommentGetListReq) (res *v1.CommentGetListRes, err error) {
	return service.Comment().GetList(ctx, req)
}
