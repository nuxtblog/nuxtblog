package history

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/history/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) HistoryClear(ctx context.Context, req *v1.HistoryClearReq) (res *v1.HistoryClearRes, err error) {
	return &v1.HistoryClearRes{}, service.History().Clear(ctx)
}
