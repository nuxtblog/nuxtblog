package plugin

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) MarketplaceSync(ctx context.Context, req *v1.MarketplaceSyncReq) (res *v1.MarketplaceSyncRes, err error) {
	return service.Plugin().SyncMarketplace(ctx)
}
