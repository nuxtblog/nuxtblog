package post

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/post/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) GetStats(ctx context.Context, req *v1.GetStatsReq) (res *v1.GetStatsRes, err error) {
	return service.Post().GetStats(ctx)
}
