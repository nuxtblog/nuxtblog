package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/follow/v1"
)

type IFollow interface {
	Follow(ctx context.Context, req *v1.UserFollowReq) (*v1.UserFollowRes, error)
	Unfollow(ctx context.Context, req *v1.UserUnfollowReq) (*v1.UserUnfollowRes, error)
	RemoveFollower(ctx context.Context, req *v1.UserRemoveFollowerReq) (*v1.UserRemoveFollowerRes, error)
	FollowStatus(ctx context.Context, req *v1.UserFollowStatusReq) (*v1.UserFollowStatusRes, error)
	Followers(ctx context.Context, req *v1.UserFollowersReq) (*v1.UserFollowersRes, error)
	Following(ctx context.Context, req *v1.UserFollowingReq) (*v1.UserFollowingRes, error)
}

var localFollow IFollow

func Follow() IFollow {
	if localFollow == nil {
		panic("implement not found for interface IFollow, forgot register?")
	}
	return localFollow
}

func RegisterFollow(i IFollow) {
	localFollow = i
}
