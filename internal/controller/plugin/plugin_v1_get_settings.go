package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginGetSettings(ctx context.Context, req *v1.PluginGetSettingsReq) (res *v1.PluginGetSettingsRes, err error) {
	return service.Plugin().GetSettings(ctx, req.Id)
}
