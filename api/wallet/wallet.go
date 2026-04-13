package wallet

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/wallet/v1"
)

type IWalletV1 interface {
	WalletBalance(ctx context.Context, req *v1.WalletBalanceReq) (res *v1.WalletBalanceRes, err error)
	WalletLedger(ctx context.Context, req *v1.WalletLedgerReq) (res *v1.WalletLedgerRes, err error)
	WalletTopup(ctx context.Context, req *v1.WalletTopupReq) (res *v1.WalletTopupRes, err error)
	WalletAdminAdjust(ctx context.Context, req *v1.WalletAdminAdjustReq) (res *v1.WalletAdminAdjustRes, err error)
	CreditsLedger(ctx context.Context, req *v1.CreditsLedgerReq) (res *v1.CreditsLedgerRes, err error)
	CreditsAdminAdjust(ctx context.Context, req *v1.CreditsAdminAdjustReq) (res *v1.CreditsAdminAdjustRes, err error)
}
