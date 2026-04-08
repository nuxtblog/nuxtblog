// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package reaction

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/reaction/v1"
)

type IReactionV1 interface {
	LikePost(ctx context.Context, req *v1.LikePostReq) (res *v1.LikePostRes, err error)
	UnlikePost(ctx context.Context, req *v1.UnlikePostReq) (res *v1.UnlikePostRes, err error)
	BookmarkPost(ctx context.Context, req *v1.BookmarkPostReq) (res *v1.BookmarkPostRes, err error)
	UnbookmarkPost(ctx context.Context, req *v1.UnbookmarkPostReq) (res *v1.UnbookmarkPostRes, err error)
	GetReaction(ctx context.Context, req *v1.GetReactionReq) (res *v1.GetReactionRes, err error)
	GetBookmarks(ctx context.Context, req *v1.GetBookmarksReq) (res *v1.GetBookmarksRes, err error)
}
