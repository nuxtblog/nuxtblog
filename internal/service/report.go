package service

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/report/v1"
)

type IReport interface {
	Create(ctx context.Context, req *v1.ReportCreateReq) (int64, error)
	List(ctx context.Context, req *v1.ReportListReq) (*v1.ReportListRes, error)
	Handle(ctx context.Context, req *v1.ReportHandleReq) error
}

var _report IReport

func Report() IReport          { return _report }
func RegisterReport(r IReport) { _report = r }
