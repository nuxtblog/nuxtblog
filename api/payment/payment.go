package payment

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"
)

type IPaymentV1 interface {
	PaymentListProviders(ctx context.Context, req *v1.PaymentListProvidersReq) (res *v1.PaymentListProvidersRes, err error)
	PaymentGetProviderConfig(ctx context.Context, req *v1.PaymentGetProviderConfigReq) (res *v1.PaymentGetProviderConfigRes, err error)
	PaymentSetProviderConfig(ctx context.Context, req *v1.PaymentSetProviderConfigReq) (res *v1.PaymentSetProviderConfigRes, err error)
}
