package service

import (
	v1 "github.com/nuxtblog/nuxtblog/api/option/v1"
	"context"
)

type (
	IOption interface {
		Get(ctx context.Context, key string) (*v1.OptionItem, error)
		GetAutoload(ctx context.Context) (map[string]string, error)
		Set(ctx context.Context, req *v1.OptionSetReq) error
		Delete(ctx context.Context, key string) error
		AdminStats(ctx context.Context) (*v1.AdminStatsRes, error)
	}
)

var (
	localOption IOption
)

func Option() IOption {
	if localOption == nil {
		panic("implement not found for interface IOption, forgot register?")
	}
	return localOption
}

func RegisterOption(i IOption) {
	localOption = i
}
