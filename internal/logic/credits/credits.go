package credits

import (
	"context"
	"fmt"

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

type sCredits struct{}

func New() service.ICredits { return &sCredits{} }

func init() {
	service.RegisterCredits(New())
}

const (
	TypeEarn        = 1
	TypeSpend       = 2
	TypeExpire      = 3
	TypeAdminAdjust = 4
)

func (s *sCredits) EnsureCredits(ctx context.Context, userID int64) error {
	cnt, _ := g.DB().Ctx(ctx).Model("user_credits").Where("user_id", userID).Count()
	if cnt > 0 {
		return nil
	}
	_, err := g.DB().Ctx(ctx).Model("user_credits").Data(g.Map{"user_id": userID}).Insert()
	return err
}

func (s *sCredits) GetLedger(ctx context.Context, req *v1.CreditsLedgerReq) (*v1.CreditsLedgerRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	m := g.DB().Ctx(ctx).Model("credits_ledger").Where("user_id", userID)
	total, _ := m.Count()

	var rows []struct {
		Id            int64  `orm:"id"`
		Type          int    `orm:"type"`
		Amount        int    `orm:"amount"`
		BalanceAfter  int    `orm:"balance_after"`
		Source        string `orm:"source"`
		ReferenceType string `orm:"reference_type"`
		ReferenceId   string `orm:"reference_id"`
		Note          string `orm:"note"`
		CreatedAt     string `orm:"created_at"`
	}
	if total > 0 {
		_ = m.Page(req.Page, req.PageSize).Order("created_at DESC").Scan(&rows)
	}

	items := make([]v1.CreditsLedgerItem, len(rows))
	for i, r := range rows {
		items[i] = v1.CreditsLedgerItem{
			ID: r.Id, Type: r.Type, Amount: r.Amount,
			BalanceAfter: r.BalanceAfter, Source: r.Source,
			ReferenceType: r.ReferenceType, ReferenceID: r.ReferenceId,
			Note: r.Note, CreatedAt: r.CreatedAt,
		}
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}
	return &v1.CreditsLedgerRes{
		Data: items, Total: total, Page: req.Page,
		PageSize: req.PageSize, TotalPages: totalPages,
	}, nil
}

func (s *sCredits) Earn(ctx context.Context, userID int64, amount int, source, refType, refID, note string) (int, error) {
	if amount <= 0 {
		return 0, gerror.NewCode(gcode.CodeInvalidParameter, "amount must be positive")
	}
	_ = s.EnsureCredits(ctx, userID)

	var balanceAfter int
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Exec(
			"UPDATE user_credits SET balance = balance + ?, total_earned = total_earned + ? WHERE user_id = ?",
			amount, amount, userID)
		if err != nil {
			return err
		}
		val, _ := tx.Model("user_credits").Where("user_id", userID).Value("balance")
		balanceAfter = val.Int()

		_, err = tx.Model("credits_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": userID, "type": TypeEarn,
			"amount": amount, "balance_after": balanceAfter,
			"source": source, "reference_type": refType,
			"reference_id": refID, "note": note,
		}).Insert()
		return err
	})
	if err != nil {
		return 0, err
	}

	event.Emit(ctx, event.CreditsEarned, payload.CreditsEarned{
		UserID: userID, Amount: amount, BalanceAfter: balanceAfter, Source: source,
	})
	return balanceAfter, nil
}

func (s *sCredits) Spend(ctx context.Context, userID int64, amount int, refType, refID, note string) (int, error) {
	if amount <= 0 {
		return 0, gerror.NewCode(gcode.CodeInvalidParameter, "amount must be positive")
	}
	_ = s.EnsureCredits(ctx, userID)

	var balanceAfter int
	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		res, err := tx.Exec(
			"UPDATE user_credits SET balance = balance - ?, total_spent = total_spent + ? WHERE user_id = ? AND balance >= ?",
			amount, amount, userID, amount)
		if err != nil {
			return err
		}
		affected, _ := res.RowsAffected()
		if affected == 0 {
			return gerror.NewCode(gcode.CodeBusinessValidationFailed, "insufficient credits")
		}

		val, _ := tx.Model("user_credits").Where("user_id", userID).Value("balance")
		balanceAfter = val.Int()

		_, err = tx.Model("credits_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": userID, "type": TypeSpend,
			"amount": -amount, "balance_after": balanceAfter,
			"source": "purchase", "reference_type": refType,
			"reference_id": refID, "note": note,
		}).Insert()
		return err
	})
	if err != nil {
		return 0, err
	}

	event.Emit(ctx, event.CreditsSpent, payload.CreditsSpent{
		UserID: userID, Amount: amount, BalanceAfter: balanceAfter,
	})
	return balanceAfter, nil
}

func (s *sCredits) AdminAdjust(ctx context.Context, req *v1.CreditsAdminAdjustReq) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}
	_ = s.EnsureCredits(ctx, req.UserID)

	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Exec(
			"UPDATE user_credits SET balance = balance + ? WHERE user_id = ?",
			req.Amount, req.UserID)
		if err != nil {
			return err
		}
		val, _ := tx.Model("user_credits").Where("user_id", req.UserID).Value("balance")
		balanceAfter := val.Int()

		_, err = tx.Model("credits_ledger").Data(g.Map{
			"id": idgen.New(), "user_id": req.UserID, "type": TypeAdminAdjust,
			"amount": req.Amount, "balance_after": balanceAfter,
			"source": "admin", "reference_type": "admin",
			"reference_id": "",
			"note": fmt.Sprintf("admin adjust: %s", req.Note),
		}).Insert()
		return err
	})
}
