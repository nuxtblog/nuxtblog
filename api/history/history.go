package history

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/history/v1"
)

type IHistory interface {
	HistoryList(ctx context.Context, req *v1.HistoryListReq) (res *v1.HistoryListRes, err error)
	HistoryClear(ctx context.Context, req *v1.HistoryClearReq) (res *v1.HistoryClearRes, err error)
}
