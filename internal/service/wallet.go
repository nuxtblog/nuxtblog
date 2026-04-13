package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/wallet/v1"
)

type IWallet interface {
	// GetBalance returns the user's wallet balance and credits.
	GetBalance(ctx context.Context, userID int64) (*v1.WalletBalanceRes, error)
	// GetLedger returns paginated wallet transaction history.
	GetLedger(ctx context.Context, req *v1.WalletLedgerReq) (*v1.WalletLedgerRes, error)
	// Topup creates a topup order for the user.
	Topup(ctx context.Context, req *v1.WalletTopupReq) (*v1.WalletTopupRes, error)
	// Spend deducts balance from the user's wallet (internal use).
	Spend(ctx context.Context, userID int64, amount int, refType, refID, note string) (balanceAfter int, err error)
	// Refund adds balance back to the user's wallet (internal use).
	Refund(ctx context.Context, userID int64, amount int, refType, refID, note string) (balanceAfter int, err error)
	// AdminAdjust adjusts a user's wallet balance (admin only).
	AdminAdjust(ctx context.Context, req *v1.WalletAdminAdjustReq) error
	// EnsureWallet creates a wallet row if it doesn't exist.
	EnsureWallet(ctx context.Context, userID int64) error
}

var _wallet IWallet

func Wallet() IWallet {
	if _wallet == nil {
		panic("implement not found for interface IWallet, forgot register?")
	}
	return _wallet
}

func RegisterWallet(i IWallet) { _wallet = i }
