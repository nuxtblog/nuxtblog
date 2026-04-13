package wallet

import (
	"context"
	"fmt"

	orderv1 "github.com/nuxtblog/nuxtblog/api/order/v1"
	v1 "github.com/nuxtblog/nuxtblog/api/wallet/v1"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sWallet struct{}

func New() service.IWallet { return &sWallet{} }

func init() {
	service.RegisterWallet(New())
}

// Ledger type constants
const (
	LedgerTopup       = 1
	LedgerSpend       = 2
	LedgerRefund      = 3
	LedgerAdminAdjust = 4
)

func (s *sWallet) EnsureWallet(ctx context.Context, userID int64) error {
	cnt, _ := g.DB().Ctx(ctx).Model("user_wallets").Where("user_id", userID).Count()
	if cnt > 0 {
		return nil
	}
	_, err := g.DB().Ctx(ctx).Model("user_wallets").Data(g.Map{"user_id": userID}).Insert()
	return err
}

func (s *sWallet) GetBalance(ctx context.Context, userID int64) (*v1.WalletBalanceRes, error) {
	_ = s.EnsureWallet(ctx, userID)

	type Row struct {
		Balance    int `orm:"balance"`
		Frozen     int `orm:"frozen"`
		TotalTopup int `orm:"total_topup"`
		TotalSpent int `orm:"total_spent"`
	}
	var row Row
	if err := g.DB().Ctx(ctx).Model("user_wallets").Where("user_id", userID).Scan(&row); err != nil {
		return nil, err
	}

	var credits int
	val, _ := g.DB().Ctx(ctx).Model("user_credits").Where("user_id", userID).Value("balance")
	if !val.IsNil() {
		credits = val.Int()
	}

	return &v1.WalletBalanceRes{
		Balance: row.Balance, Frozen: row.Frozen,
		TotalTopup: row.TotalTopup, TotalSpent: row.TotalSpent,
		Credits: credits,
	}, nil
}

func (s *sWallet) GetLedger(ctx context.Context, req *v1.WalletLedgerReq) (*v1.WalletLedgerRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	m := g.DB().Ctx(ctx).Model("wallet_ledger").Where("user_id", userID)
	total, _ := m.Count()

	var rows []struct {
		Id            int64  `orm:"id"`
		Type          int    `orm:"type"`
		Amount        int    `orm:"amount"`
		BalanceAfter  int    `orm:"balance_after"`
		ReferenceType string `orm:"reference_type"`
		ReferenceId   string `orm:"reference_id"`
		Note          string `orm:"note"`
		CreatedAt     string `orm:"created_at"`
	}
	if total > 0 {
		_ = m.Page(req.Page, req.PageSize).Order("created_at DESC").Scan(&rows)
	}

	items := make([]v1.LedgerItem, len(rows))
	for i, r := range rows {
		items[i] = v1.LedgerItem{
			ID: r.Id, Type: r.Type, Amount: r.Amount,
			BalanceAfter: r.BalanceAfter, ReferenceType: r.ReferenceType,
			ReferenceID: r.ReferenceId, Note: r.Note, CreatedAt: r.CreatedAt,
		}
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}
	return &v1.WalletLedgerRes{
		Data: items, Total: total, Page: req.Page,
		PageSize: req.PageSize, TotalPages: totalPages,
	}, nil
}

func (s *sWallet) Topup(ctx context.Context, req *v1.WalletTopupReq) (*v1.WalletTopupRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	orderRes, err := service.Order().Create(ctx, &orderv1.OrderCreateReq{
		Items: []orderv1.OrderItem{{
			ItemType:  orderv1.ItemTypeTopup,
			ItemID:    "topup",
			Title:     "Wallet Topup",
			UnitPrice: req.Amount,
			Quantity:  1,
		}},
	})
	if err != nil {
		return nil, err
	}
	return &v1.WalletTopupRes{OrderID: orderRes.ID, OrderNo: orderRes.OrderNo}, nil
}

// CreditTopup directly credits the wallet (called after topup order is paid).
func CreditTopup(ctx context.Context, userID int64, amount int, refType, refID string) (int, error) {
	if amount <= 0 {
		return 0, gerror.NewCode(gcode.CodeInvalidParameter, "amount must be positive")
	}
	_ = service.Wallet().EnsureWallet(ctx, userID)

	var balanceAfter int
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Exec(
			"UPDATE user_wallets SET balance = balance + ?, total_topup = total_topup + ? WHERE user_id = ?",
			amount, amount, userID)
		if err != nil {
			return err
		}
		val, _ := tx.Model("user_wallets").Where("user_id", userID).Value("balance")
		balanceAfter = val.Int()

		_, err = tx.Model("wallet_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": userID, "type": LedgerTopup,
			"amount": amount, "balance_after": balanceAfter,
			"reference_type": refType, "reference_id": refID, "note": "topup",
		}).Insert()
		return err
	})
	if err != nil {
		return 0, err
	}

	event.Emit(ctx, event.WalletTopup, payload.WalletTopup{
		UserID: userID, Amount: amount, BalanceAfter: balanceAfter,
	})
	return balanceAfter, nil
}

func (s *sWallet) Spend(ctx context.Context, userID int64, amount int, refType, refID, note string) (int, error) {
	if amount <= 0 {
		return 0, gerror.NewCode(gcode.CodeInvalidParameter, "amount must be positive")
	}
	_ = s.EnsureWallet(ctx, userID)

	var balanceAfter int
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Exec(
			"UPDATE user_wallets SET balance = balance - ?, total_spent = total_spent + ? WHERE user_id = ? AND balance >= ?",
			amount, amount, userID, amount)
		if err != nil {
			return err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "insufficient balance")
		}

		val, _ := tx.Model("user_wallets").Where("user_id", userID).Value("balance")
		balanceAfter = val.Int()

		_, err = tx.Model("wallet_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": userID, "type": LedgerSpend,
			"amount": -amount, "balance_after": balanceAfter,
			"reference_type": refType, "reference_id": refID, "note": note,
		}).Insert()
		return err
	})
	if err != nil {
		return 0, err
	}

	event.Emit(ctx, event.WalletSpend, payload.WalletSpend{
		UserID: userID, Amount: amount, BalanceAfter: balanceAfter,
	})
	return balanceAfter, nil
}

func (s *sWallet) Refund(ctx context.Context, userID int64, amount int, refType, refID, note string) (int, error) {
	if amount <= 0 {
		return 0, gerror.NewCode(gcode.CodeInvalidParameter, "amount must be positive")
	}
	_ = s.EnsureWallet(ctx, userID)

	var balanceAfter int
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Exec(
			"UPDATE user_wallets SET balance = balance + ? WHERE user_id = ?",
			amount, userID)
		if err != nil {
			return err
		}
		val, _ := tx.Model("user_wallets").Where("user_id", userID).Value("balance")
		balanceAfter = val.Int()

		_, err = tx.Model("wallet_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": userID, "type": LedgerRefund,
			"amount": amount, "balance_after": balanceAfter,
			"reference_type": refType, "reference_id": refID, "note": note,
		}).Insert()
		return err
	})
	return balanceAfter, err
}

func (s *sWallet) AdminAdjust(ctx context.Context, req *v1.WalletAdminAdjustReq) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}
	_ = s.EnsureWallet(ctx, req.UserID)

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Exec(
			"UPDATE user_wallets SET balance = balance + ? WHERE user_id = ?",
			req.Amount, req.UserID)
		if err != nil {
			return err
		}
		val, _ := tx.Model("user_wallets").Where("user_id", req.UserID).Value("balance")
		balanceAfter := val.Int()

		_, err = tx.Model("wallet_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": req.UserID, "type": LedgerAdminAdjust,
			"amount": req.Amount, "balance_after": balanceAfter,
			"reference_type": "admin", "reference_id": "",
			"note": fmt.Sprintf("admin adjust: %s", req.Note),
		}).Insert()
		return err
	})
}
