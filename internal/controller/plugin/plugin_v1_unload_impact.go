package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginUnloadImpact(ctx context.Context, req *v1.PluginUnloadImpactReq) (res *v1.PluginUnloadImpactRes, err error) {
	return service.Plugin().UnloadImpact(ctx, req.Id)
}
