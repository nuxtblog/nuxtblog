package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginUninstall(ctx context.Context, req *v1.PluginUninstallReq) (res *v1.PluginUninstallRes, err error) {
	return &v1.PluginUninstallRes{}, service.Plugin().Uninstall(ctx, req.Id)
}
