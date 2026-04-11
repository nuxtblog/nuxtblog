package system

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/system/v1"
)

type ISystemV1 interface {
	Info(ctx context.Context, req *v1.InfoReq) (res *v1.InfoRes, err error)
}
