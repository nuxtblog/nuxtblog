package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/friendlink/v1"
)

type IFriendlink interface {
	Create(ctx context.Context, name, url, logo, description string, sortOrder, status int) (int64, error)
	ListAdmin(ctx context.Context, page, size int) ([]*v1.FriendlinkItem, int, error)
	Update(ctx context.Context, id int64, name, url, logo, description string, sortOrder, status int) error
	Delete(ctx context.Context, id int64) error
	ListPublic(ctx context.Context) ([]*v1.FriendlinkItem, error)
}

var localFriendlink IFriendlink

func Friendlink() IFriendlink {
	if localFriendlink == nil {
		panic("implement not found for interface IFriendlink, forgot register?")
	}
	return localFriendlink
}

func RegisterFriendlink(i IFriendlink) {
	localFriendlink = i
}
