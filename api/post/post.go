// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package post

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/post/v1"
)

type IPostV1 interface {
	PostCreate(ctx context.Context, req *v1.PostCreateReq) (res *v1.PostCreateRes, err error)
	PostDelete(ctx context.Context, req *v1.PostDeleteReq) (res *v1.PostDeleteRes, err error)
	PostUpdate(ctx context.Context, req *v1.PostUpdateReq) (res *v1.PostUpdateRes, err error)
	PostGetOne(ctx context.Context, req *v1.PostGetOneReq) (res *v1.PostGetOneRes, err error)
	PostGetBySlug(ctx context.Context, req *v1.PostGetBySlugReq) (res *v1.PostGetBySlugRes, err error)
	PostGetList(ctx context.Context, req *v1.PostGetListReq) (res *v1.PostGetListRes, err error)
	PostSeoUpdate(ctx context.Context, req *v1.PostSeoUpdateReq) (res *v1.PostSeoUpdateRes, err error)
	PostRevisionList(ctx context.Context, req *v1.PostRevisionListReq) (res *v1.PostRevisionListRes, err error)
	PostRevisionRestore(ctx context.Context, req *v1.PostRevisionRestoreReq) (res *v1.PostRevisionRestoreRes, err error)
	GetStats(ctx context.Context, req *v1.GetStatsReq) (res *v1.GetStatsRes, err error)
	PostTrash(ctx context.Context, req *v1.PostTrashReq) (res *v1.PostTrashRes, err error)
	PostRestore(ctx context.Context, req *v1.PostRestoreReq) (res *v1.PostRestoreRes, err error)
	PostBatch(ctx context.Context, req *v1.PostBatchReq) (res *v1.PostBatchRes, err error)
	PostMetaUpdate(ctx context.Context, req *v1.PostMetaUpdateReq) (res *v1.PostMetaUpdateRes, err error)
	PostMetaGet(ctx context.Context, req *v1.PostMetaGetReq) (res *v1.PostMetaGetRes, err error)
}
