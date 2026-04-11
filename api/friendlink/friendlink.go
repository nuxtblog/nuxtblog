package friendlink

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/friendlink/v1"
)

type IFriendlinkAdmin interface {
	FriendlinkCreate(ctx context.Context, req *v1.FriendlinkCreateReq) (res *v1.FriendlinkCreateRes, err error)
	FriendlinkAdminList(ctx context.Context, req *v1.FriendlinkAdminListReq) (res *v1.FriendlinkAdminListRes, err error)
	FriendlinkUpdate(ctx context.Context, req *v1.FriendlinkUpdateReq) (res *v1.FriendlinkUpdateRes, err error)
	FriendlinkDelete(ctx context.Context, req *v1.FriendlinkDeleteReq) (res *v1.FriendlinkDeleteRes, err error)
}

type IFriendlinkPublic interface {
	FriendlinkList(ctx context.Context, req *v1.FriendlinkListReq) (res *v1.FriendlinkListRes, err error)
}
