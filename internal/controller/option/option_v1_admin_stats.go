package option

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) AdminStats(ctx context.Context, req *v1.AdminStatsReq) (res *v1.AdminStatsRes, err error) {
	return service.Option().AdminStats(ctx)
}
