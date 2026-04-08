package service

import "context"

type ICheckin interface {
	DoCheckin(ctx context.Context, userID int64) (alreadyCheckedIn bool, streak int, err error)
	GetStatus(ctx context.Context, userID int64) (checkedInToday bool, streak int, err error)
}

var localCheckin ICheckin

func Checkin() ICheckin {
	if localCheckin == nil {
		panic("implement not found for interface ICheckin, forgot register?")
	}
	return localCheckin
}

func RegisterCheckin(i ICheckin) {
	localCheckin = i
}
