package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MarketplaceList(ctx context.Context, req *v1.MarketplaceListReq) (res *v1.MarketplaceListRes, err error) {
	return service.Plugin().GetMarketplace(ctx, req.Keyword, req.Type)
}
