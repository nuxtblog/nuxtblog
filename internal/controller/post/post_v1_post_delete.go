package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostDelete(ctx context.Context, req *v1.PostDeleteReq) (res *v1.PostDeleteRes, err error) {
	return &v1.PostDeleteRes{}, service.Post().Delete(ctx, req.Id)
}
