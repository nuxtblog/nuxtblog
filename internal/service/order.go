package service

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/order/v1"
)

type IOrder interface {
	// Create creates a new order.
	Create(ctx context.Context, req *v1.OrderCreateReq) (*v1.OrderCreateRes, error)
	// GetList returns paginated order list for a user or admin.
	GetList(ctx context.Context, req *v1.OrderListReq) (*v1.OrderListRes, error)
	// GetByID returns order detail by ID.
	GetByID(ctx context.Context, id int64) (*v1.OrderDetailRes, error)
	// Pay initiates payment for an order, returns payment info.
	Pay(ctx context.Context, req *v1.OrderPayReq) (*v1.OrderPayRes, error)
	// Cancel cancels a pending order.
	Cancel(ctx context.Context, id int64) error
	// HandleNotify processes payment provider callback.
	HandleNotify(ctx context.Context, provider string, body []byte, headers map[string]string) error
	// Refund refunds an order (admin only).
	Refund(ctx context.Context, req *v1.OrderRefundReq) error
	// Complete marks a paid order as completed (delivers items).
	Complete(ctx context.Context, orderID int64) error
	// CheckPurchase checks if a user has purchased a specific item.
	CheckPurchase(ctx context.Context, req *v1.PurchaseCheckReq) (*v1.PurchaseCheckRes, error)
	// GetRevenueStats returns revenue statistics (admin only).
	GetRevenueStats(ctx context.Context) (*v1.RevenueStatsRes, error)
}

var _order IOrder

func Order() IOrder {
	if _order == nil {
		panic("implement not found for interface IOrder, forgot register?")
	}
	return _order
}

func RegisterOrder(i IOrder) { _order = i }
