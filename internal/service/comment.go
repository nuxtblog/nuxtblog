package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/comment/v1"
	"context"
)

type (
	IComment interface {
		AdminGetList(ctx context.Context, req *v1.CommentAdminGetListReq) (*v1.CommentAdminGetListRes, error)
		Create(ctx context.Context, req *v1.CommentCreateReq) (*v1.CommentCreateRes, error)
		Delete(ctx context.Context, id int64) error
		GetList(ctx context.Context, req *v1.CommentGetListReq) (*v1.CommentGetListRes, error)
		UpdateStatus(ctx context.Context, id int64, status v1.CommentReviewStatus) error
		UpdateContent(ctx context.Context, id int64, content string) error
	}
)

var (
	localComment IComment
)

func Comment() IComment {
	if localComment == nil {
		panic("implement not found for interface IComment, forgot register?")
	}
	return localComment
}

func RegisterComment(i IComment) {
	localComment = i
}
