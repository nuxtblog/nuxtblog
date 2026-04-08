package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostCreate(ctx context.Context, req *v1.PostCreateReq) (res *v1.PostCreateRes, err error) {
	id, err := service.Post().Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.PostCreateRes{Id: id}, nil
}
