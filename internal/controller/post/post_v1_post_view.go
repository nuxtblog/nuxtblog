package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostView(ctx context.Context, req *v1.PostViewReq) (res *v1.PostViewRes, err error) {
	if err = service.Post().IncrementView(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.PostViewRes{}, nil
}
