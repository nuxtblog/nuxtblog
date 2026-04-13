package order

import (
	"context"
	"fmt"
	"strconv"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/order/v1"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	plugin "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sOrder struct{}

func New() service.IOrder { return &sOrder{} }

func init() {
	service.RegisterOrder(New())
}

// generateOrderNo creates a unique order number based on timestamp + random suffix.
func generateOrderNo() string {
	return fmt.Sprintf("%s%d", time.Now().Format("20060102150405"), idgen.New()%100000)
}

// ── Create ──────────────────────────────────────────────────────────────────

func (s *sOrder) Create(ctx context.Context, req *v1.OrderCreateReq) (*v1.OrderCreateRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	if len(req.Items) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "items required")
	}

	// Calculate total
	var totalAmount int
	for _, item := range req.Items {
		totalAmount += item.UnitPrice * item.Quantity
	}

	// Run filter:price.calculate to allow membership discounts etc.
	if filtered, err := plugin.Filter(ctx, plugin.FilterPriceCalculate, map[string]any{
		"user_id":      userID,
		"total_amount": totalAmount,
		"items":        req.Items,
	}); err == nil {
		if v, ok := filtered["total_amount"].(int); ok {
			totalAmount = v
		}
	}

	// Run filter:order.create
	if _, err := plugin.Filter(ctx, plugin.FilterOrderCreate, map[string]any{
		"user_id":      userID,
		"total_amount": totalAmount,
		"items":        req.Items,
	}); err != nil {
		return nil, err
	}

	orderID := idgen.New()
	orderNo := generateOrderNo()

	// Apply balance and credits if requested
	balanceUsed := req.UseBalance
	creditsUsed := req.UseCredits
	paidAmount := totalAmount

	// Deduct credits first
	if creditsUsed > 0 && creditsUsed <= paidAmount {
		paidAmount -= creditsUsed
	} else {
		creditsUsed = 0
	}

	// Deduct balance
	if balanceUsed > 0 && balanceUsed <= paidAmount {
		paidAmount -= balanceUsed
	} else if balanceUsed > paidAmount {
		balanceUsed = paidAmount
		paidAmount = 0
	}

	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	// Set order expiry (30 minutes for unpaid orders)
	expiresAt := time.Now().Add(30 * time.Minute)

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		_, err := tx.Model("orders").Data(g.Map{
			"id": orderID, "order_no": orderNo, "user_id": userID,
			"status": v1.OrderStatusPending, "total_amount": totalAmount,
			"paid_amount": 0, "credits_used": creditsUsed,
			"balance_used": balanceUsed, "currency": currency,
			"expires_at": expiresAt,
		}).Insert()
		if err != nil {
			return err
		}

		for _, item := range req.Items {
			_, err = tx.Model("order_items").Data(g.Map{
				"id": idgen.New(), "order_id": orderID,
				"item_type": item.ItemType, "item_id": item.ItemID,
				"title": item.Title, "unit_price": item.UnitPrice,
				"quantity": item.Quantity, "snapshot": item.Snapshot,
			}).Insert()
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	event.Emit(ctx, event.OrderCreated, payload.OrderCreated{
		OrderID: orderID, OrderNo: orderNo, UserID: userID,
		Amount: totalAmount, ItemType: req.Items[0].ItemType,
		ItemID: req.Items[0].ItemID,
	})

	return &v1.OrderCreateRes{
		ID: orderID, OrderNo: orderNo,
		Status: v1.OrderStatusPending, Amount: totalAmount,
	}, nil
}

// ── GetList ─────────────────────────────────────────────────────────────────

func (s *sOrder) GetList(ctx context.Context, req *v1.OrderListReq) (*v1.OrderListRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	m := g.DB().Ctx(ctx).Model("orders").Where("user_id", userID)
	if req.Status != nil {
		m = m.Where("status", *req.Status)
	}

	total, _ := m.Count()
	var orders []OrderRow
	if total > 0 {
		_ = m.Page(req.Page, req.PageSize).Order("created_at DESC").Scan(&orders)
	}

	items := make([]v1.OrderListItem, len(orders))
	for i, o := range orders {
		items[i] = o.toListItem()
		items[i].Items = s.getOrderItems(ctx, o.Id)
	}

	totalPages := 0
	if req.PageSize > 0 {
		totalPages = (total + req.PageSize - 1) / req.PageSize
	}
	return &v1.OrderListRes{
		Data: items, Total: total, Page: req.Page,
		PageSize: req.PageSize, TotalPages: totalPages,
	}, nil
}

// ── GetByID ─────────────────────────────────────────────────────────────────

func (s *sOrder) GetByID(ctx context.Context, id int64) (*v1.OrderDetailRes, error) {
	var o OrderRow
	if err := g.DB().Ctx(ctx).Model("orders").Where("id", id).Scan(&o); err != nil {
		return nil, err
	}
	if o.Id == 0 {
		return nil, gerror.NewCode(gcode.CodeNotFound, "order not found")
	}

	item := o.toListItem()
	item.Items = s.getOrderItems(ctx, id)
	return &v1.OrderDetailRes{OrderListItem: item}, nil
}

// ── Pay ─────────────────────────────────────────────────────────────────────

func (s *sOrder) Pay(ctx context.Context, req *v1.OrderPayReq) (*v1.OrderPayRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	var o OrderRow
	if err := g.DB().Ctx(ctx).Model("orders").Where("id", req.ID).Scan(&o); err != nil {
		return nil, err
	}
	if o.Id == 0 || o.UserId != userID {
		return nil, gerror.NewCode(gcode.CodeNotFound, "order not found")
	}
	if o.Status != v1.OrderStatusPending {
		return nil, gerror.NewCode(gcode.CodeBusinessValidationFailed, "order is not pending")
	}

	remainingAmount := o.TotalAmount - o.CreditsUsed - o.BalanceUsed

	// Balance payment: instant deduction
	if req.Provider == "balance" {
		// Deduct credits first if specified
		if o.CreditsUsed > 0 {
			_, err := service.Credits().Spend(ctx, userID, o.CreditsUsed,
				"order", strconv.FormatInt(o.Id, 10), "order payment")
			if err != nil {
				return nil, err
			}
		}

		// Deduct remaining from wallet
		payAmount := remainingAmount
		if o.BalanceUsed > 0 {
			payAmount = o.BalanceUsed
		}
		if payAmount <= 0 {
			payAmount = remainingAmount
		}

		_, err := service.Wallet().Spend(ctx, userID, payAmount,
			"order", strconv.FormatInt(o.Id, 10), "order payment")
		if err != nil {
			return nil, err
		}

		// Record transaction
		_, err = g.DB().Ctx(ctx).Model("transactions").Data(g.Map{
			"id": idgen.New(), "order_id": o.Id, "user_id": userID,
			"type": 1, "provider": "balance", "provider_tx_id": "",
			"amount": payAmount, "status": 2,
		}).Insert()
		if err != nil {
			return nil, err
		}

		// Mark order as paid
		_, err = g.DB().Ctx(ctx).Model("orders").Where("id", o.Id).Data(g.Map{
			"status": v1.OrderStatusPaid, "paid_amount": payAmount,
		}).Update()
		if err != nil {
			return nil, err
		}

		event.Emit(ctx, event.OrderPaid, payload.OrderPaid{
			OrderID: o.Id, OrderNo: o.OrderNo, UserID: userID,
			Amount: payAmount, Provider: "balance",
		})

		// Auto-complete the order
		_ = s.Complete(ctx, o.Id)

		return &v1.OrderPayRes{Paid: true}, nil
	}

	// External payment providers (alipay, paypal) — placeholder
	// In production, this would call the payment provider's CreatePayment API
	// and return a payment URL or QR code.
	return &v1.OrderPayRes{
		PaymentURL: fmt.Sprintf("/pay/%s?order=%s&amount=%d", req.Provider, o.OrderNo, remainingAmount),
	}, nil
}

// ── Cancel ──────────────────────────────────────────────────────────────────

func (s *sOrder) Cancel(ctx context.Context, id int64) error {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return gerror.NewCode(gcode.CodeNotAuthorized, "login required")
	}

	var o OrderRow
	if err := g.DB().Ctx(ctx).Model("orders").Where("id", id).Scan(&o); err != nil {
		return err
	}
	if o.Id == 0 || o.UserId != userID {
		return gerror.NewCode(gcode.CodeNotFound, "order not found")
	}
	if o.Status != v1.OrderStatusPending {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "only pending orders can be cancelled")
	}

	_, err := g.DB().Ctx(ctx).Model("orders").Where("id", id).Data(g.Map{
		"status": v1.OrderStatusCancelled,
	}).Update()
	if err != nil {
		return err
	}

	event.Emit(ctx, event.OrderCancelled, payload.OrderCancelled{
		OrderID: o.Id, OrderNo: o.OrderNo, UserID: userID,
	})
	return nil
}

// ── HandleNotify ────────────────────────────────────────────────────────────

func (s *sOrder) HandleNotify(ctx context.Context, provider string, body []byte, headers map[string]string) error {
	// TODO: Implement actual provider-specific verification (Alipay RSA2, PayPal webhook)
	// For now this is a placeholder. Real implementation needs:
	// 1. Verify signature
	// 2. Parse order_no and amount
	// 3. Find order, mark as paid
	// 4. Complete order
	g.Log().Infof(ctx, "payment notify received: provider=%s", provider)
	return nil
}

// ── Refund ──────────────────────────────────────────────────────────────────

func (s *sOrder) Refund(ctx context.Context, req *v1.OrderRefundReq) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}

	var o OrderRow
	if err := g.DB().Ctx(ctx).Model("orders").Where("id", req.ID).Scan(&o); err != nil {
		return err
	}
	if o.Id == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "order not found")
	}
	if o.Status != v1.OrderStatusPaid && o.Status != v1.OrderStatusCompleted {
		return gerror.NewCode(gcode.CodeBusinessValidationFailed, "order cannot be refunded")
	}

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// Refund balance
		if o.PaidAmount > 0 {
			_, err := service.Wallet().Refund(ctx, o.UserId, o.PaidAmount,
				"order", strconv.FormatInt(o.Id, 10), "order refund")
			if err != nil {
				return err
			}
		}

		// Refund credits
		if o.CreditsUsed > 0 {
			_, err := service.Credits().Earn(ctx, o.UserId, o.CreditsUsed,
				"refund", "order", strconv.FormatInt(o.Id, 10), "order refund")
			if err != nil {
				return err
			}
		}

		// Update order status
		_, err := tx.Model("orders").Where("id", o.Id).Data(g.Map{
			"status": v1.OrderStatusRefunded,
		}).Update()
		if err != nil {
			return err
		}

		// Record refund transaction
		_, err = tx.Model("transactions").Data(g.Map{
			"id": idgen.New(), "order_id": o.Id, "user_id": o.UserId,
			"type": 2, "provider": "balance", "amount": o.PaidAmount,
			"status": 2,
		}).Insert()
		if err != nil {
			return err
		}

		// Remove purchase records
		_, err = tx.Model("user_purchases").Where("order_id", o.Id).Delete()
		return err
	})
	if err != nil {
		return err
	}

	event.Emit(ctx, event.OrderRefunded, payload.OrderRefunded{
		OrderID: o.Id, OrderNo: o.OrderNo, UserID: o.UserId, Amount: o.PaidAmount,
	})
	return nil
}

// ── Complete ────────────────────────────────────────────────────────────────

func (s *sOrder) Complete(ctx context.Context, orderID int64) error {
	var o OrderRow
	if err := g.DB().Ctx(ctx).Model("orders").Where("id", orderID).Scan(&o); err != nil {
		return err
	}
	if o.Id == 0 {
		return gerror.NewCode(gcode.CodeNotFound, "order not found")
	}

	// Get order items
	items := s.getOrderItems(ctx, orderID)

	err := g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		// Mark order as completed
		_, err := tx.Model("orders").Where("id", orderID).Data(g.Map{
			"status": v1.OrderStatusCompleted,
		}).Update()
		if err != nil {
			return err
		}

		// Deliver items
		for _, item := range items {
			switch item.ItemType {
			case v1.ItemTypePostUnlock, v1.ItemTypeDownload, v1.ItemTypeProduct:
				// Record purchase
				_, _ = tx.Model("user_purchases").Data(g.Map{
					"user_id": o.UserId, "object_type": item.ItemType,
					"object_id": item.ItemID, "order_id": orderID,
				}).Insert()

				event.Emit(ctx, event.ContentUnlocked, payload.ContentUnlocked{
					UserID: o.UserId, ObjectType: item.ItemType,
					ObjectID: item.ItemID, OrderID: orderID,
				})

			case v1.ItemTypeMembership:
				// Activate membership
				tierID, _ := strconv.ParseInt(item.ItemID, 10, 64)
				if tierID > 0 {
					_ = service.Membership().Activate(ctx, o.UserId, tierID, orderID)
				}

			case v1.ItemTypeTopup:
				// Topup is handled by the commerce listener on order.completed event
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(items) > 0 {
		event.Emit(ctx, event.OrderCompleted, payload.OrderCompleted{
			OrderID: o.Id, OrderNo: o.OrderNo, UserID: o.UserId,
			ItemType: items[0].ItemType, ItemID: items[0].ItemID,
		})
	}
	return nil
}

// ── CheckPurchase ───────────────────────────────────────────────────────────

func (s *sOrder) CheckPurchase(ctx context.Context, req *v1.PurchaseCheckReq) (*v1.PurchaseCheckRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return &v1.PurchaseCheckRes{Purchased: false}, nil
	}

	type Row struct {
		OrderId int64 `orm:"order_id"`
	}
	var row Row
	err := g.DB().Ctx(ctx).Model("user_purchases").
		Where("user_id", userID).
		Where("object_type", req.ObjectType).
		Where("object_id", req.ObjectID).
		Scan(&row)
	if err != nil || row.OrderId == 0 {
		// Also check membership access
		hasAccess, _, err := service.Membership().CheckAccess(ctx, userID)
		if err == nil && hasAccess {
			return &v1.PurchaseCheckRes{Purchased: true}, nil
		}
		return &v1.PurchaseCheckRes{Purchased: false}, nil
	}

	return &v1.PurchaseCheckRes{Purchased: true, OrderID: &row.OrderId}, nil
}

// ── GetRevenueStats ─────────────────────────────────────────────────────────

func (s *sOrder) GetRevenueStats(ctx context.Context) (*v1.RevenueStatsRes, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}

	db := g.DB().Ctx(ctx)

	// Total revenue (paid + completed orders)
	totalRev, _ := db.Model("orders").
		WhereIn("status", g.Slice{v1.OrderStatusPaid, v1.OrderStatusCompleted}).
		Sum("paid_amount")

	// Today's revenue
	today := time.Now().Format("2006-01-02")
	todayRev, _ := db.Model("orders").
		WhereIn("status", g.Slice{v1.OrderStatusPaid, v1.OrderStatusCompleted}).
		WhereLike("created_at", today+"%").
		Sum("paid_amount")

	totalOrders, _ := db.Model("orders").Count()
	todayOrders, _ := db.Model("orders").WhereLike("created_at", today+"%").Count()
	activeMembers, _ := db.Model("user_memberships").Where("status", 1).Count()
	pendingOrders, _ := db.Model("orders").Where("status", v1.OrderStatusPending).Count()

	return &v1.RevenueStatsRes{
		TotalRevenue:  int(totalRev),
		TodayRevenue:  int(todayRev),
		TotalOrders:   totalOrders,
		TodayOrders:   todayOrders,
		ActiveMembers: activeMembers,
		PendingOrders: pendingOrders,
	}, nil
}

// ── Helpers ─────────────────────────────────────────────────────────────────

type OrderRow struct {
	Id          int64  `orm:"id"`
	OrderNo     string `orm:"order_no"`
	UserId      int64  `orm:"user_id"`
	Status      int    `orm:"status"`
	TotalAmount int    `orm:"total_amount"`
	PaidAmount  int    `orm:"paid_amount"`
	CreditsUsed int    `orm:"credits_used"`
	BalanceUsed int    `orm:"balance_used"`
	Currency    string `orm:"currency"`
	CreatedAt   string `orm:"created_at"`
	UpdatedAt   string `orm:"updated_at"`
}

func (o OrderRow) toListItem() v1.OrderListItem {
	return v1.OrderListItem{
		ID: o.Id, OrderNo: o.OrderNo, UserID: o.UserId,
		Status: o.Status, TotalAmount: o.TotalAmount,
		PaidAmount: o.PaidAmount, CreditsUsed: o.CreditsUsed,
		BalanceUsed: o.BalanceUsed, Currency: o.Currency,
		CreatedAt: o.CreatedAt, UpdatedAt: o.UpdatedAt,
	}
}

func (s *sOrder) getOrderItems(ctx context.Context, orderID int64) []v1.OrderItem {
	var rows []struct {
		ItemType  string `orm:"item_type"`
		ItemId    string `orm:"item_id"`
		Title     string `orm:"title"`
		UnitPrice int    `orm:"unit_price"`
		Quantity  int    `orm:"quantity"`
		Snapshot  string `orm:"snapshot"`
	}
	_ = g.DB().Ctx(ctx).Model("order_items").Where("order_id", orderID).Scan(&rows)

	items := make([]v1.OrderItem, len(rows))
	for i, r := range rows {
		items[i] = v1.OrderItem{
			ItemType: r.ItemType, ItemID: r.ItemId, Title: r.Title,
			UnitPrice: r.UnitPrice, Quantity: r.Quantity, Snapshot: r.Snapshot,
		}
	}
	return items
}
