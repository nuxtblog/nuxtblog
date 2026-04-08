package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostRestore(ctx context.Context, req *v1.PostRestoreReq) (res *v1.PostRestoreRes, err error) {
	return &v1.PostRestoreRes{}, service.Post().Restore(ctx, req.Id)
}
