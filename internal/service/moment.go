package service

import (
	momentv1 "github.com/nuxtblog/nuxtblog/api/moment/v1"
	"context"
)

type IMoment interface {
	Create(ctx context.Context, req *momentv1.MomentCreateReq) (int64, error)
	Update(ctx context.Context, req *momentv1.MomentUpdateReq) error
	Delete(ctx context.Context, id int64) error
	GetById(ctx context.Context, id int64) (*momentv1.MomentItem, error)
	GetList(ctx context.Context, req *momentv1.MomentGetListReq) (*momentv1.MomentGetListRes, error)
	IncrementView(ctx context.Context, id int64) error
}

var localMoment IMoment

func Moment() IMoment {
	if localMoment == nil {
		panic("implement not found for interface IMoment, forgot register?")
	}
	return localMoment
}

func RegisterMoment(i IMoment) {
	localMoment = i
}
