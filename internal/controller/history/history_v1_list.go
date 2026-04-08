package history

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/history/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) HistoryList(ctx context.Context, req *v1.HistoryListReq) (res *v1.HistoryListRes, err error) {
	return service.History().List(ctx, req)
}
