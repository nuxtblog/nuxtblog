package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostBatchUpdate(ctx context.Context, req *v1.PostBatchUpdateReq) (res *v1.PostBatchUpdateRes, err error) {
	affected, err := service.Post().BatchUpdate(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.PostBatchUpdateRes{Affected: affected}, nil
}
