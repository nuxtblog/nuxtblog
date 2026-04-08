package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginBatchUpdate(ctx context.Context, req *v1.PluginBatchUpdateReq) (*v1.PluginBatchUpdateRes, error) {
	return service.Plugin().BatchUpdate(ctx, req.Ids)
}
