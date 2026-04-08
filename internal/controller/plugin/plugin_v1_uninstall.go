package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginUninstall(ctx context.Context, req *v1.PluginUninstallReq) (res *v1.PluginUninstallRes, err error) {
	needRestart, err := service.Plugin().Uninstall(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.PluginUninstallRes{NeedRestart: needRestart}, nil
}
