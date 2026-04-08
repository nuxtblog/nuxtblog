package report

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/report/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ReportList(ctx context.Context, req *v1.ReportListReq) (res *v1.ReportListRes, err error) {
	return service.Report().List(ctx, req)
}
