package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/history/v1"
)

type IHistory interface {
	List(ctx context.Context, req *v1.HistoryListReq) (*v1.HistoryListRes, error)
	Clear(ctx context.Context) error
}

var _history IHistory

func History() IHistory          { return _history }
func RegisterHistory(h IHistory) { _history = h }
