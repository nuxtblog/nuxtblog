// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package comment

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/comment/v1"
)

type ICommentV1 interface {
	CommentCreate(ctx context.Context, req *v1.CommentCreateReq) (res *v1.CommentCreateRes, err error)
	CommentDelete(ctx context.Context, req *v1.CommentDeleteReq) (res *v1.CommentDeleteRes, err error)
	CommentUpdateStatus(ctx context.Context, req *v1.CommentUpdateStatusReq) (res *v1.CommentUpdateStatusRes, err error)
	CommentUpdateContent(ctx context.Context, req *v1.CommentUpdateContentReq) (res *v1.CommentUpdateContentRes, err error)
	CommentGetList(ctx context.Context, req *v1.CommentGetListReq) (res *v1.CommentGetListRes, err error)
	CommentAdminGetList(ctx context.Context, req *v1.CommentAdminGetListReq) (res *v1.CommentAdminGetListRes, err error)
}
