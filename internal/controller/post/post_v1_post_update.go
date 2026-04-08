package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostUpdate(ctx context.Context, req *v1.PostUpdateReq) (res *v1.PostUpdateRes, err error) {
	return &v1.PostUpdateRes{}, service.Post().Update(ctx, req)
}
