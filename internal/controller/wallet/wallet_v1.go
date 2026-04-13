package wallet

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/wallet/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

func (c *ControllerV1) WalletBalance(ctx context.Context, req *v1.WalletBalanceReq) (*v1.WalletBalanceRes, error) {
	userID, _ := middleware.GetCurrentUserID(ctx)
	return service.Wallet().GetBalance(ctx, userID)
}

func (c *ControllerV1) WalletLedger(ctx context.Context, req *v1.WalletLedgerReq) (*v1.WalletLedgerRes, error) {
	return service.Wallet().GetLedger(ctx, req)
}

func (c *ControllerV1) WalletTopup(ctx context.Context, req *v1.WalletTopupReq) (*v1.WalletTopupRes, error) {
	return service.Wallet().Topup(ctx, req)
}

func (c *ControllerV1) WalletAdminAdjust(ctx context.Context, req *v1.WalletAdminAdjustReq) (*v1.WalletAdminAdjustRes, error) {
	return &v1.WalletAdminAdjustRes{}, service.Wallet().AdminAdjust(ctx, req)
}

func (c *ControllerV1) CreditsLedger(ctx context.Context, req *v1.CreditsLedgerReq) (*v1.CreditsLedgerRes, error) {
	return service.Credits().GetLedger(ctx, req)
}

func (c *ControllerV1) CreditsAdminAdjust(ctx context.Context, req *v1.CreditsAdminAdjustReq) (*v1.CreditsAdminAdjustRes, error) {
	return &v1.CreditsAdminAdjustRes{}, service.Credits().AdminAdjust(ctx, req)
}
