package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginUpdateSettings(ctx context.Context, req *v1.PluginUpdateSettingsReq) (res *v1.PluginUpdateSettingsRes, err error) {
	return &v1.PluginUpdateSettingsRes{}, service.Plugin().UpdateSettings(ctx, req.Id, req.Values)
}
