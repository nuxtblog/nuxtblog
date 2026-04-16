package service

import (
	"context"
	v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"
)

// CreatePaymentReq is the internal request type for creating a payment.
type CreatePaymentReq struct {
	OrderNo   string
	Amount    int
	Currency  string
	Subject   string
	NotifyURL string
	ReturnURL string
	ClientIP  string
}

// CreatePaymentRes is the internal response from creating a payment.
type CreatePaymentRes struct {
	PaymentURL string
	QRCode     string
	Method     string
}

// NotifyResult is the internal parsed result from a payment provider callback.
type NotifyResult struct {
	Success      bool
	OrderNo      string
	Amount       int
	ProviderTxID string
	RawResponse  string
}

type IPayment interface {
	ListProviders(ctx context.Context) (*v1.PaymentListProvidersRes, error)
	GetProviderConfig(ctx context.Context, slug string) (*v1.PaymentGetProviderConfigRes, error)
	SetProviderConfig(ctx context.Context, slug string, config map[string]interface{}) (*v1.PaymentSetProviderConfigRes, error)
	// ListEnabledProviders returns basic info for enabled providers (no admin required).
	ListEnabledProviders(ctx context.Context) ([]v1.ProviderBasicInfo, error)
	// CreatePayment creates a payment via the specified provider.
	CreatePayment(ctx context.Context, provider string, req CreatePaymentReq) (*CreatePaymentRes, error)
	// VerifyNotify verifies a payment provider callback.
	VerifyNotify(ctx context.Context, provider string, body []byte, headers map[string]string) (*NotifyResult, error)
}

var _payment IPayment

func Payment() IPayment          { return _payment }
func RegisterPayment(p IPayment) { _payment = p }
