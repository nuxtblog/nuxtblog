package service

import (
	reactionv1 "github.com/nuxtblog/nuxtblog/api/reaction/v1"
	"context"
)

type IReaction interface {
	Like(ctx context.Context, postId, userId int64) error
	Unlike(ctx context.Context, postId, userId int64) error
	Bookmark(ctx context.Context, postId, userId int64) error
	Unbookmark(ctx context.Context, postId, userId int64) error
	GetStatus(ctx context.Context, postId, userId int64) (*reactionv1.GetReactionRes, error)
	GetBookmarks(ctx context.Context, userId int64, page, size int) (*reactionv1.GetBookmarksRes, error)
}

var localReaction IReaction

func Reaction() IReaction {
	if localReaction == nil {
		panic("implement not found for interface IReaction, forgot register?")
	}
	return localReaction
}

func RegisterReaction(i IReaction) {
	localReaction = i
}
