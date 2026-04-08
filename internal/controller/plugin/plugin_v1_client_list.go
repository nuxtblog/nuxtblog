package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginClientList(ctx context.Context, req *v1.PluginClientListReq) (res *v1.PluginClientListRes, err error) {
	return service.Plugin().ClientList(ctx)
}
