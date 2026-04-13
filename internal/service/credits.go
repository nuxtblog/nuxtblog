package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/wallet/v1"
)

type ICredits interface {
	// GetLedger returns paginated credits transaction history.
	GetLedger(ctx context.Context, req *v1.CreditsLedgerReq) (*v1.CreditsLedgerRes, error)
	// Earn adds credits to the user's account.
	Earn(ctx context.Context, userID int64, amount int, source, refType, refID, note string) (balanceAfter int, err error)
	// Spend deducts credits from the user's account.
	Spend(ctx context.Context, userID int64, amount int, refType, refID, note string) (balanceAfter int, err error)
	// AdminAdjust adjusts a user's credits (admin only).
	AdminAdjust(ctx context.Context, req *v1.CreditsAdminAdjustReq) error
	// EnsureCredits creates a credits row if it doesn't exist.
	EnsureCredits(ctx context.Context, userID int64) error
}

var _credits ICredits

func Credits() ICredits {
	if _credits == nil {
		panic("implement not found for interface ICredits, forgot register?")
	}
	return _credits
}

func RegisterCredits(i ICredits) { _credits = i }
