package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginGetErrors(ctx context.Context, req *v1.PluginGetErrorsReq) (res *v1.PluginGetErrorsRes, err error) {
	return service.Plugin().GetErrors(ctx, req.Id)
}
