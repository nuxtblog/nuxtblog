package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginUpdate(ctx context.Context, req *v1.PluginUpdateReq) (res *v1.PluginUpdateRes, err error) {
	item, err := service.Plugin().Update(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.PluginUpdateRes{Item: *item}, nil
}
