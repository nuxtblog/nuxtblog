package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostSeoUpdate(ctx context.Context, req *v1.PostSeoUpdateReq) (res *v1.PostSeoUpdateRes, err error) {
	return &v1.PostSeoUpdateRes{}, service.Post().SeoUpdate(ctx, req)
}
