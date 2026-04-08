package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostTrash(ctx context.Context, req *v1.PostTrashReq) (res *v1.PostTrashRes, err error) {
	return &v1.PostTrashRes{}, service.Post().Trash(ctx, req.Id)
}
