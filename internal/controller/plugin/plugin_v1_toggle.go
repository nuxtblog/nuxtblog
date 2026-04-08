package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginToggle(ctx context.Context, req *v1.PluginToggleReq) (res *v1.PluginToggleRes, err error) {
	return &v1.PluginToggleRes{}, service.Plugin().Toggle(ctx, req.Id, req.Enabled)
}
