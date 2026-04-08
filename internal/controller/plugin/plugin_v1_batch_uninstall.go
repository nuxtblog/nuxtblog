package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginBatchUninstall(ctx context.Context, req *v1.PluginBatchUninstallReq) (res *v1.PluginBatchUninstallRes, err error) {
	res = &v1.PluginBatchUninstallRes{}
	for _, id := range req.Ids {
		needRestart, e := service.Plugin().Uninstall(ctx, id)
		if e != nil {
			res.Failed = append(res.Failed, id)
			continue
		}
		res.Succeeded = append(res.Succeeded, id)
		if needRestart {
			res.NeedRestart = true
		}
	}
	return res, nil
}
