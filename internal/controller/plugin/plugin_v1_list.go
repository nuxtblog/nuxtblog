package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginList(ctx context.Context, req *v1.PluginListReq) (res *v1.PluginListRes, err error) {
	return service.Plugin().List(ctx)
}
