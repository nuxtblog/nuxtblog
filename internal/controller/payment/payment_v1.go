package payment

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) PaymentListProviders(ctx context.Context, req *v1.PaymentListProvidersReq) (*v1.PaymentListProvidersRes, error) {
	return service.Payment().ListProviders(ctx)
}

func (c *ControllerV1) PaymentGetProviderConfig(ctx context.Context, req *v1.PaymentGetProviderConfigReq) (*v1.PaymentGetProviderConfigRes, error) {
	return service.Payment().GetProviderConfig(ctx, req.Slug)
}

func (c *ControllerV1) PaymentSetProviderConfig(ctx context.Context, req *v1.PaymentSetProviderConfigReq) (*v1.PaymentSetProviderConfigRes, error) {
	return service.Payment().SetProviderConfig(ctx, req.Slug, req.Config)
}
