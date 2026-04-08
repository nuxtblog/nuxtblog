package report

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/report/v1"
)

type IReport interface {
	ReportCreate(ctx context.Context, req *v1.ReportCreateReq) (res *v1.ReportCreateRes, err error)
	ReportList(ctx context.Context, req *v1.ReportListReq) (res *v1.ReportListRes, err error)
	ReportHandle(ctx context.Context, req *v1.ReportHandleReq) (res *v1.ReportHandleRes, err error)
}
