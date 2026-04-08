package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginPreview(ctx context.Context, req *v1.PluginPreviewReq) (res *v1.PluginPreviewRes, err error) {
	return service.Plugin().Preview(ctx, req.Repo)
}
