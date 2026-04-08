// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package media

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/media/v1"
)

type IMediaV1 interface {
	MediaUpload(ctx context.Context, req *v1.MediaUploadReq) (res *v1.MediaUploadRes, err error)
	MediaDelete(ctx context.Context, req *v1.MediaDeleteReq) (res *v1.MediaDeleteRes, err error)
	MediaUpdate(ctx context.Context, req *v1.MediaUpdateReq) (res *v1.MediaUpdateRes, err error)
	MediaGetOne(ctx context.Context, req *v1.MediaGetOneReq) (res *v1.MediaGetOneRes, err error)
	MediaGetList(ctx context.Context, req *v1.MediaGetListReq) (res *v1.MediaGetListRes, err error)
	MediaGetStats(ctx context.Context, req *v1.MediaGetStatsReq) (res *v1.MediaGetStatsRes, err error)
	MediaCategoryList(ctx context.Context, req *v1.MediaCategoryListReq) (res *v1.MediaCategoryListRes, err error)
	MediaCategoryUpdate(ctx context.Context, req *v1.MediaCategoryUpdateReq) (res *v1.MediaCategoryUpdateRes, err error)
}
