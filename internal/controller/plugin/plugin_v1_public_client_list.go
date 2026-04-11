package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginPublicClientList(ctx context.Context, req *v1.PluginPublicClientReq) (res *v1.PluginPublicClientRes, err error) {
	return service.Plugin().PublicClientList(ctx)
}
