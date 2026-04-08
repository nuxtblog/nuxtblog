package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostVerifyPassword(ctx context.Context, req *v1.PostVerifyPasswordReq) (res *v1.PostVerifyPasswordRes, err error) {
	valid, err := service.Post().VerifyPassword(ctx, req.Id, req.Password)
	if err != nil {
		return nil, err
	}
	return &v1.PostVerifyPasswordRes{Valid: valid}, nil
}
