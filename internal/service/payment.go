package service

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"
)

type IPayment interface {
	ListProviders(ctx context.Context) (*v1.PaymentListProvidersRes, error)
	GetProviderConfig(ctx context.Context, slug string) (*v1.PaymentGetProviderConfigRes, error)
	SetProviderConfig(ctx context.Context, slug string, config map[string]interface{}) (*v1.PaymentSetProviderConfigRes, error)
}

var _payment IPayment

func Payment() IPayment          { return _payment }
func RegisterPayment(p IPayment) { _payment = p }
