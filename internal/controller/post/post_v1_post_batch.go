package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PostBatch(ctx context.Context, req *v1.PostBatchReq) (res *v1.PostBatchRes, err error) {
	affected, err := service.Post().Batch(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.PostBatchRes{Affected: affected}, nil
}
