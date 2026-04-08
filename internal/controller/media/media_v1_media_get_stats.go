package media

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MediaGetStats(ctx context.Context, req *v1.MediaGetStatsReq) (res *v1.MediaGetStatsRes, err error) {
	return service.Media().GetStats(ctx)
}
