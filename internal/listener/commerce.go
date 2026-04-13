package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	walletlogic "github.com/nuxtblog/nuxtblog/internal/logic/wallet"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
)

func registerCommerceListeners() {
	// Award credits on check-in
	event.OnAsync(event.CheckinDone, onCheckinAwardCredits)
	// Award credits on comment
	event.OnAsync(event.CommentCreated, onCommentAwardCredits)
	// Handle topup order completion → credit wallet
	event.OnAsync(event.OrderCompleted, onOrderCompletedTopup)
	// Membership expire check (also triggered by cron, but we listen for
	// membership.expired to notify the user)
	event.OnAsync(event.MembershipExpired, onMembershipExpired)
}

// ── Check-in → credits ──────────────────────────────────────────────────────

func onCheckinAwardCredits(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.CheckinDone)
	if !ok || p.AlreadyCheckedIn {
		return nil
	}

	amount := getCreditsRule(ctx, "checkin", 10)
	if amount <= 0 {
		return nil
	}

	_, err := service.Credits().Earn(ctx, p.UserID, amount, "checkin", "checkin", "", "daily check-in")
	if err != nil {
		g.Log().Warningf(ctx, "award checkin credits user=%d: %v", p.UserID, err)
	}
	return nil
}

// ── Comment → credits ───────────────────────────────────────────────────────

func onCommentAwardCredits(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.CommentCreated)
	if !ok || p.AuthorID == 0 {
		return nil
	}

	amount := getCreditsRule(ctx, "comment", 5)
	dailyLimit := getCreditsRule(ctx, "comment_daily_limit", 10)
	if amount <= 0 {
		return nil
	}

	// Check daily comment credits count
	todayCount, err := g.DB().Ctx(ctx).Model("credits_ledger").
		Where("user_id", p.AuthorID).
		Where("source", "comment").
		WhereGTE("created_at", g.DB().Raw("date('now')")).
		Count()
	if err != nil {
		g.Log().Warningf(ctx, "check comment credits limit: %v", err)
		return nil
	}
	if todayCount >= dailyLimit {
		return nil
	}

	_, err = service.Credits().Earn(ctx, p.AuthorID, amount, "comment", "comment",
		fmt.Sprintf("%d", p.CommentID), "comment reward")
	if err != nil {
		g.Log().Warningf(ctx, "award comment credits user=%d: %v", p.AuthorID, err)
	}
	return nil
}

// ── Order completed → topup wallet ──────────────────────────────────────────

func onOrderCompletedTopup(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.OrderCompleted)
	if !ok || p.ItemType != "topup" {
		return nil
	}

	// Get order total to know how much to credit
	type Row struct {
		TotalAmount int `orm:"total_amount"`
	}
	var row Row
	err := g.DB().Ctx(ctx).Model("orders").Where("id", p.OrderID).Scan(&row)
	if err != nil || row.TotalAmount == 0 {
		g.Log().Warningf(ctx, "topup order %d: cannot read amount: %v", p.OrderID, err)
		return nil
	}

	_, err = walletlogic.CreditTopup(ctx, p.UserID, row.TotalAmount, "order",
		strconv.FormatInt(p.OrderID, 10))
	if err != nil {
		g.Log().Warningf(ctx, "credit topup wallet user=%d order=%d: %v", p.UserID, p.OrderID, err)
	}
	return nil
}

// ── Membership expired → notify user ────────────────────────────────────────

func onMembershipExpired(ctx context.Context, e event.Event) error {
	p, ok := e.Payload.(payload.MembershipExpired)
	if !ok {
		return nil
	}
	_ = service.Notification().Create(ctx,
		"system", "membership_expired",
		nil, "", "",
		p.UserID,
		"membership", nil,
		"Membership expired", "/user/membership",
		g.I18n().Tf(ctx, "notification.membership_expired"),
	)
	return nil
}

// ── helpers ─────────────────────────────────────────────────────────────────

func getCreditsRule(ctx context.Context, key string, defaultVal int) int {
	val, err := g.DB().Ctx(ctx).Model("options").
		Where("key", "credits_rules").
		Value("value")
	if err != nil || val.IsNil() {
		return defaultVal
	}
	var rules map[string]int
	if err := json.Unmarshal([]byte(val.String()), &rules); err != nil {
		return defaultVal
	}
	if v, ok := rules[key]; ok {
		return v
	}
	return defaultVal
}
