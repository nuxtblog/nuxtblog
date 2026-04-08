package moment

import (
	"context"
	"github.com/nuxtblog/nuxtblog/api/moment/v1"
)

type IMomentV1 interface {
	MomentCreate(ctx context.Context, req *v1.MomentCreateReq) (res *v1.MomentCreateRes, err error)
	MomentUpdate(ctx context.Context, req *v1.MomentUpdateReq) (res *v1.MomentUpdateRes, err error)
	MomentDelete(ctx context.Context, req *v1.MomentDeleteReq) (res *v1.MomentDeleteRes, err error)
	MomentGetOne(ctx context.Context, req *v1.MomentGetOneReq) (res *v1.MomentGetOneRes, err error)
	MomentGetList(ctx context.Context, req *v1.MomentGetListReq) (res *v1.MomentGetListRes, err error)
	MomentView(ctx context.Context, req *v1.MomentViewReq) (res *v1.MomentViewRes, err error)
}
