package plugin

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PluginStyles(ctx context.Context, req *v1.PluginStylesReq) (res *v1.PluginStylesRes, err error) {
	css, err := service.Plugin().GetStyles(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.PluginStylesRes{CSS: css}, nil
}
