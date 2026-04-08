package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginBatchUninstall(ctx context.Context, req *v1.PluginBatchUninstallReq) (res *v1.PluginBatchUninstallRes, err error) {
	res = &v1.PluginBatchUninstallRes{}
	for _, id := range req.Ids {
		if e := service.Plugin().Uninstall(ctx, id); e != nil {
			res.Failed = append(res.Failed, id)
		} else {
			res.Succeeded = append(res.Succeeded, id)
		}
	}
	return res, nil
}
