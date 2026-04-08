// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package option

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/option/v1"
)

type IOptionV1 interface {
	OptionGet(ctx context.Context, req *v1.OptionGetReq) (res *v1.OptionGetRes, err error)
	OptionGetAutoload(ctx context.Context, req *v1.OptionGetAutoloadReq) (res *v1.OptionGetAutoloadRes, err error)
	OptionSet(ctx context.Context, req *v1.OptionSetReq) (res *v1.OptionSetRes, err error)
	OptionDelete(ctx context.Context, req *v1.OptionDeleteReq) (res *v1.OptionDeleteRes, err error)
	AdminStats(ctx context.Context, req *v1.AdminStatsReq) (res *v1.AdminStatsRes, err error)
	StorageBackends(ctx context.Context, req *v1.StorageBackendsReq) (res *v1.StorageBackendsRes, err error)
}
