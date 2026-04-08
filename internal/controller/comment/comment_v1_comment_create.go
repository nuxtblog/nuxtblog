package comment

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) CommentCreate(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error) {
	return service.Comment().Create(ctx, req)
}
