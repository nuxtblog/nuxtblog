package follow

import (
	"context"

	"github.com/nuxtblog/nuxtblog/api/follow/v1"
)

// IFollowPublicV1 — routes registered in the public group (no auth required).
type IFollowPublicV1 interface {
	UserFollowers(ctx context.Context, req *v1.UserFollowersReq) (res *v1.UserFollowersRes, err error)
	UserFollowing(ctx context.Context, req *v1.UserFollowingReq) (res *v1.UserFollowingRes, err error)
}

// IFollowAuthV1 — routes registered in the authenticated group (JWT required).
type IFollowAuthV1 interface {
	UserFollow(ctx context.Context, req *v1.UserFollowReq) (res *v1.UserFollowRes, err error)
	UserUnfollow(ctx context.Context, req *v1.UserUnfollowReq) (res *v1.UserUnfollowRes, err error)
	UserRemoveFollower(ctx context.Context, req *v1.UserRemoveFollowerReq) (res *v1.UserRemoveFollowerRes, err error)
	UserFollowStatus(ctx context.Context, req *v1.UserFollowStatusReq) (res *v1.UserFollowStatusRes, err error)
}
