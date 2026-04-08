package report

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/report/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) ReportCreate(ctx context.Context, req *v1.ReportCreateReq) (res *v1.ReportCreateRes, err error) {
	id, err := service.Report().Create(ctx, req)
	return &v1.ReportCreateRes{Id: id}, err
}
