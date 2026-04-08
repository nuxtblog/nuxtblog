package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginInstall(ctx context.Context, req *v1.PluginInstallReq) (res *v1.PluginInstallRes, err error) {
	item, err := service.Plugin().Install(ctx, req.RepoUrl, req.ExpectedVersion)
	if err != nil {
		return nil, err
	}
	return &v1.PluginInstallRes{Item: *item}, nil
}
