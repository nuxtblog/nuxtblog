package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginGetStats(ctx context.Context, req *v1.PluginGetStatsReq) (res *v1.PluginGetStatsRes, err error) {
	return service.Plugin().GetStats(ctx, req.Id)
}
