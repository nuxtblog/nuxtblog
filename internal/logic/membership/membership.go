package membership

import (
	"context"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/order/v1"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sMembership struct{}

func New() service.IMembership { return &sMembership{} }

func init() {
	service.RegisterMembership(New())
}

// ── ListTiers (public) ──────────────────────────────────────────────────────

func (s *sMembership) ListTiers(ctx context.Context) (*v1.MembershipTierListRes, error) {
	var rows []tierRow
	err := g.DB().Ctx(ctx).Model("membership_tiers").
		Where("status", 1).
		Order("sort_order ASC, id ASC").
		Scan(&rows)
	if err != nil {
		return nil, err
	}

	items := make([]v1.MembershipTierItem, len(rows))
	for i, r := range rows {
		items[i] = r.toItem()
	}
	return &v1.MembershipTierListRes{Items: items}, nil
}

// ── GetUserMembership ───────────────────────────────────────────────────────

func (s *sMembership) GetUserMembership(ctx context.Context, userID int64) (*v1.UserMembershipRes, error) {
	type Row struct {
		TierId    int64  `orm:"tier_id"`
		Status    int    `orm:"status"`
		ExpiresAt string `orm:"expires_at"`
		AutoRenew int    `orm:"auto_renew"`
	}
	var row Row
	err := g.DB().Ctx(ctx).Model("user_memberships").
		Where("user_id", userID).
		Where("status", 1).
		Order("expires_at DESC").
		Scan(&row)
	if err != nil || row.TierId == 0 {
		return &v1.UserMembershipRes{Active: false}, nil
	}

	var tier tierRow
	_ = g.DB().Ctx(ctx).Model("membership_tiers").Where("id", row.TierId).Scan(&tier)

	item := tier.toItem()
	return &v1.UserMembershipRes{
		Active:    true,
		Tier:      &item,
		ExpiresAt: row.ExpiresAt,
		AutoRenew: row.AutoRenew == 1,
	}, nil
}

// ── Activate ────────────────────────────────────────────────────────────────

func (s *sMembership) Activate(ctx context.Context, userID int64, tierID int64, orderID int64) error {
	var tier tierRow
	if err := g.DB().Ctx(ctx).Model("membership_tiers").Where("id", tierID).Scan(&tier); err != nil {
		return err
	}
	if tier.Id == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "tier not found")
	}

	expiresAt := time.Now().AddDate(0, 0, tier.DurationDays)

	// Check existing active membership
	type ExistRow struct {
		Id        int64  `orm:"id"`
		ExpiresAt string `orm:"expires_at"`
	}
	var existing ExistRow
	_ = g.DB().Ctx(ctx).Model("user_memberships").
		Where("user_id", userID).
		Where("status", 1).
		Scan(&existing)

	if existing.Id > 0 {
		// Extend existing membership
		if t, err := time.Parse(time.RFC3339, existing.ExpiresAt); err == nil && t.After(time.Now()) {
			expiresAt = t.AddDate(0, 0, tier.DurationDays)
		}
		_, err := g.DB().Ctx(ctx).Model("user_memberships").Where("id", existing.Id).Data(g.Map{
			"tier_id": tierID, "expires_at": expiresAt, "order_id": orderID,
		}).Update()
		if err != nil {
			return err
		}

		event.Emit(ctx, event.MembershipRenewed, payload.MembershipRenewed{
			UserID: userID, TierID: tierID, TierSlug: tier.Slug, OrderID: orderID,
		})
	} else {
		// Create new membership
		_, err := g.DB().Ctx(ctx).Model("user_memberships").Data(g.Map{
			"id": idgen.New(), "user_id": userID, "tier_id": tierID,
			"status": 1, "expires_at": expiresAt,
			"auto_renew": 0, "order_id": orderID,
		}).Insert()
		if err != nil {
			return err
		}

		event.Emit(ctx, event.MembershipActivated, payload.MembershipActivated{
			UserID: userID, TierID: tierID, TierSlug: tier.Slug, OrderID: orderID,
		})
	}
	return nil
}

// ── CheckAccess ─────────────────────────────────────────────────────────────

func (s *sMembership) CheckAccess(ctx context.Context, userID int64) (bool, int, error) {
	type Row struct {
		TierId int64 `orm:"tier_id"`
	}
	var row Row
	err := g.DB().Ctx(ctx).Model("user_memberships").
		Where("user_id", userID).
		Where("status", 1).
		WhereGTE("expires_at", time.Now()).
		Scan(&row)
	if err != nil || row.TierId == 0 {
		return false, 0, nil
	}

	type TierRow struct {
		AccessAll   int `orm:"access_all"`
		DiscountPct int `orm:"discount_pct"`
	}
	var tier TierRow
	_ = g.DB().Ctx(ctx).Model("membership_tiers").Where("id", row.TierId).Scan(&tier)

	hasAccess := tier.AccessAll == 1
	return hasAccess, tier.DiscountPct, nil
}

// ── ExpireCheck (cron) ──────────────────────────────────────────────────────

func (s *sMembership) ExpireCheck(ctx context.Context) {
	type Row struct {
		Id     int64  `orm:"id"`
		UserId int64  `orm:"user_id"`
		TierId int64  `orm:"tier_id"`
	}
	var rows []Row
	_ = g.DB().Ctx(ctx).Model("user_memberships").
		Where("status", 1).
		WhereLT("expires_at", time.Now()).
		Scan(&rows)

	for _, r := range rows {
		_, _ = g.DB().Ctx(ctx).Model("user_memberships").Where("id", r.Id).Data(g.Map{
			"status": 2, // expired
		}).Update()

		var slug string
		val, _ := g.DB().Ctx(ctx).Model("membership_tiers").Where("id", r.TierId).Value("slug")
		slug = val.String()

		event.Emit(ctx, event.MembershipExpired, payload.MembershipExpired{
			UserID: r.UserId, TierID: r.TierId, TierSlug: slug,
		})
	}
	if len(rows) > 0 {
		g.Log().Infof(ctx, "membership expire check: %d memberships expired", len(rows))
	}
}

// ── Admin CRUD ──────────────────────────────────────────────────────────────

func (s *sMembership) AdminCreateTier(ctx context.Context, req *v1.MembershipTierCreateReq) (*v1.MembershipTierItem, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}

	id := idgen.New()
	features := req.Features
	if features == "" {
		features = "[]"
	}
	_, err := g.DB().Ctx(ctx).Model("membership_tiers").Data(g.Map{
		"id": id, "name": req.Name, "slug": req.Slug,
		"description": req.Description, "price": req.Price,
		"duration_days": req.DurationDays, "discount_pct": req.DiscountPct,
		"access_all": boolToInt(req.AccessAll), "credits_monthly": req.CreditsMonthly,
		"features": features, "status": req.Status, "sort_order": req.SortOrder,
	}).Insert()
	if err != nil {
		return nil, err
	}

	return &v1.MembershipTierItem{
		ID: id, Name: req.Name, Slug: req.Slug,
		Description: req.Description, Price: req.Price,
		DurationDays: req.DurationDays, DiscountPct: req.DiscountPct,
		AccessAll: req.AccessAll, CreditsMonthly: req.CreditsMonthly,
		Features: features, Status: req.Status, SortOrder: req.SortOrder,
	}, nil
}

func (s *sMembership) AdminUpdateTier(ctx context.Context, req *v1.MembershipTierUpdateReq) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}

	data := g.Map{}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.Description != "" {
		data["description"] = req.Description
	}
	if req.Price != nil {
		data["price"] = *req.Price
	}
	if req.DurationDays != nil {
		data["duration_days"] = *req.DurationDays
	}
	if req.DiscountPct != nil {
		data["discount_pct"] = *req.DiscountPct
	}
	if req.AccessAll != nil {
		data["access_all"] = boolToInt(*req.AccessAll)
	}
	if req.CreditsMonthly != nil {
		data["credits_monthly"] = *req.CreditsMonthly
	}
	if req.Features != "" {
		data["features"] = req.Features
	}
	if req.Status != nil {
		data["status"] = *req.Status
	}
	if req.SortOrder != nil {
		data["sort_order"] = *req.SortOrder
	}
	if len(data) == 0 {
		return nil
	}
	_, err := g.DB().Ctx(ctx).Model("membership_tiers").Where("id", req.ID).Data(data).Update()
	return err
}

func (s *sMembership) AdminDeleteTier(ctx context.Context, id int64) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}
	_, err := g.DB().Ctx(ctx).Model("membership_tiers").Where("id", id).Delete()
	return err
}

func (s *sMembership) AdminListSubscribers(ctx context.Context, req *v1.MembershipSubscriberListReq) (*v1.MembershipSubscriberListRes, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}

	m := g.DB().Ctx(ctx).Model("user_memberships um").
		LeftJoin("users u", "u.id = um.user_id").
		LeftJoin("membership_tiers mt", "mt.id = um.tier_id").
		Fields("um.user_id, u.username, mt.name as tier_name, um.status, um.started_at, um.expires_at, um.auto_renew")

	total, _ := m.Count()
	var rows []struct {
		UserId    int64  `orm:"user_id"`
		Username  string `orm:"username"`
		TierName  string `orm:"tier_name"`
		Status    int    `orm:"status"`
		StartedAt string `orm:"started_at"`
		ExpiresAt string `orm:"expires_at"`
		AutoRenew int    `orm:"auto_renew"`
	}
	if total > 0 {
		_ = m.Page(req.Page, req.PageSize).Order("um.created_at DESC").Scan(&rows)
	}

	items := make([]v1.SubscriberItem, len(rows))
	for i, r := range rows {
		items[i] = v1.SubscriberItem{
			UserID: r.UserId, Username: r.Username, TierName: r.TierName,
			Status: r.Status, StartedAt: r.StartedAt, ExpiresAt: r.ExpiresAt,
			AutoRenew: r.AutoRenew == 1,
		}
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}
	return &v1.MembershipSubscriberListRes{
		Data: items, Total: total, Page: req.Page,
		PageSize: req.PageSize, TotalPages: totalPages,
	}, nil
}

// ── helpers ─────────────────────────────────────────────────────────────────

type tierRow struct {
	Id             int64  `orm:"id"`
	Name           string `orm:"name"`
	Slug           string `orm:"slug"`
	Description    string `orm:"description"`
	Price          int    `orm:"price"`
	DurationDays   int    `orm:"duration_days"`
	DiscountPct    int    `orm:"discount_pct"`
	AccessAll      int    `orm:"access_all"`
	CreditsMonthly int    `orm:"credits_monthly"`
	Features       string `orm:"features"`
	Status         int    `orm:"status"`
	SortOrder      int    `orm:"sort_order"`
}

func (r tierRow) toItem() v1.MembershipTierItem {
	return v1.MembershipTierItem{
		ID: r.Id, Name: r.Name, Slug: r.Slug,
		Description: r.Description, Price: r.Price,
		DurationDays: r.DurationDays, DiscountPct: r.DiscountPct,
		AccessAll: r.AccessAll == 1, CreditsMonthly: r.CreditsMonthly,
		Features: r.Features, Status: r.Status, SortOrder: r.SortOrder,
	}
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
