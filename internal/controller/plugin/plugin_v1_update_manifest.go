package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginUpdateManifest(ctx context.Context, req *v1.PluginUpdateManifestReq) (res *v1.PluginUpdateManifestRes, err error) {
	return nil, service.Plugin().UpdateManifest(ctx, req.Id, req.Manifest)
}
