package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginGetManifest(ctx context.Context, req *v1.PluginGetManifestReq) (res *v1.PluginGetManifestRes, err error) {
	return service.Plugin().GetManifest(ctx, req.Id)
}
