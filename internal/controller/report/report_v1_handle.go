package report

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/report/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ReportHandle(ctx context.Context, req *v1.ReportHandleReq) (res *v1.ReportHandleRes, err error) {
	return &v1.ReportHandleRes{}, service.Report().Handle(ctx, req)
}
