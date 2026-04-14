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
	MediaLink(ctx context.Context, req *v1.MediaLinkReq) (res *v1.MediaLinkRes, err error)
	MediaLocalize(ctx context.Context, req *v1.MediaLocalizeReq) (res *v1.MediaLocalizeRes, err error)
	ExtensionGroupList(ctx context.Context, req *v1.ExtensionGroupListReq) (res *v1.ExtensionGroupListRes, err error)
	ExtensionGroupSave(ctx context.Context, req *v1.ExtensionGroupSaveReq) (res *v1.ExtensionGroupSaveRes, err error)
	FormatPolicyList(ctx context.Context, req *v1.FormatPolicyListReq) (res *v1.FormatPolicyListRes, err error)
	FormatPolicyCreate(ctx context.Context, req *v1.FormatPolicyCreateReq) (res *v1.FormatPolicyCreateRes, err error)
	FormatPolicyUpdate(ctx context.Context, req *v1.FormatPolicyUpdateReq) (res *v1.FormatPolicyUpdateRes, err error)
	FormatPolicyDelete(ctx context.Context, req *v1.FormatPolicyDeleteReq) (res *v1.FormatPolicyDeleteRes, err error)
}
