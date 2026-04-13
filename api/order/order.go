package order

import (
	"context"

	v1 "github.com/nuxtblog/nuxtblog/api/order/v1"
)

type IOrderV1 interface {
	OrderCreate(ctx context.Context, req *v1.OrderCreateReq) (res *v1.OrderCreateRes, err error)
	OrderList(ctx context.Context, req *v1.OrderListReq) (res *v1.OrderListRes, err error)
	OrderDetail(ctx context.Context, req *v1.OrderDetailReq) (res *v1.OrderDetailRes, err error)
	OrderPay(ctx context.Context, req *v1.OrderPayReq) (res *v1.OrderPayRes, err error)
	OrderCancel(ctx context.Context, req *v1.OrderCancelReq) (res *v1.OrderCancelRes, err error)
	PurchaseCheck(ctx context.Context, req *v1.PurchaseCheckReq) (res *v1.PurchaseCheckRes, err error)
	PaymentNotify(ctx context.Context, req *v1.PaymentNotifyReq) (res *v1.PaymentNotifyRes, err error)
	AdminOrderList(ctx context.Context, req *v1.AdminOrderListReq) (res *v1.AdminOrderListRes, err error)
	OrderRefund(ctx context.Context, req *v1.OrderRefundReq) (res *v1.OrderRefundRes, err error)
	RevenueStats(ctx context.Context, req *v1.RevenueStatsReq) (res *v1.RevenueStatsRes, err error)
	MembershipTierList(ctx context.Context, req *v1.MembershipTierListReq) (res *v1.MembershipTierListRes, err error)
	MembershipTierCreate(ctx context.Context, req *v1.MembershipTierCreateReq) (res *v1.MembershipTierCreateRes, err error)
	MembershipTierUpdate(ctx context.Context, req *v1.MembershipTierUpdateReq) (res *v1.MembershipTierUpdateRes, err error)
	MembershipTierDelete(ctx context.Context, req *v1.MembershipTierDeleteReq) (res *v1.MembershipTierDeleteRes, err error)
	UserMembership(ctx context.Context, req *v1.UserMembershipReq) (res *v1.UserMembershipRes, err error)
	MembershipSubscriberList(ctx context.Context, req *v1.MembershipSubscriberListReq) (res *v1.MembershipSubscriberListRes, err error)
}
