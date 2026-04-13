package order

import (
	"context"
	"io"
	"net/http"

	v1 "github.com/nuxtblog/nuxtblog/api/order/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/net/ghttp"
)

func (c *ControllerV1) OrderCreate(ctx context.Context, req *v1.OrderCreateReq) (*v1.OrderCreateRes, error) {
	return service.Order().Create(ctx, req)
}

func (c *ControllerV1) OrderList(ctx context.Context, req *v1.OrderListReq) (*v1.OrderListRes, error) {
	return service.Order().GetList(ctx, req)
}

func (c *ControllerV1) OrderDetail(ctx context.Context, req *v1.OrderDetailReq) (*v1.OrderDetailRes, error) {
	return service.Order().GetByID(ctx, req.ID)
}

func (c *ControllerV1) OrderPay(ctx context.Context, req *v1.OrderPayReq) (*v1.OrderPayRes, error) {
	return service.Order().Pay(ctx, req)
}

func (c *ControllerV1) OrderCancel(ctx context.Context, req *v1.OrderCancelReq) (*v1.OrderCancelRes, error) {
	return &v1.OrderCancelRes{}, service.Order().Cancel(ctx, req.ID)
}

func (c *ControllerV1) PurchaseCheck(ctx context.Context, req *v1.PurchaseCheckReq) (*v1.PurchaseCheckRes, error) {
	return service.Order().CheckPurchase(ctx, req)
}

func (c *ControllerV1) PaymentNotify(ctx context.Context, req *v1.PaymentNotifyReq) (*v1.PaymentNotifyRes, error) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return &v1.PaymentNotifyRes{}, nil
	}
	body, _ := io.ReadAll(r.Request.Body)
	headers := make(map[string]string)
	for k := range r.Request.Header {
		headers[http.CanonicalHeaderKey(k)] = r.Request.Header.Get(k)
	}
	return &v1.PaymentNotifyRes{}, service.Order().HandleNotify(ctx, req.Provider, body, headers)
}

func (c *ControllerV1) AdminOrderList(ctx context.Context, req *v1.AdminOrderListReq) (*v1.AdminOrderListRes, error) {
	res, err := service.Order().GetList(ctx, &v1.OrderListReq{
		Page: req.Page, PageSize: req.PageSize, Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.AdminOrderListRes{
		Data: res.Data, Total: res.Total, Page: res.Page,
		PageSize: res.PageSize, TotalPages: res.TotalPages,
	}, nil
}

func (c *ControllerV1) OrderRefund(ctx context.Context, req *v1.OrderRefundReq) (*v1.OrderRefundRes, error) {
	return &v1.OrderRefundRes{}, service.Order().Refund(ctx, req)
}

func (c *ControllerV1) RevenueStats(ctx context.Context, req *v1.RevenueStatsReq) (*v1.RevenueStatsRes, error) {
	return service.Order().GetRevenueStats(ctx)
}

func (c *ControllerV1) MembershipTierList(ctx context.Context, req *v1.MembershipTierListReq) (*v1.MembershipTierListRes, error) {
	return service.Membership().ListTiers(ctx)
}

func (c *ControllerV1) MembershipTierCreate(ctx context.Context, req *v1.MembershipTierCreateReq) (*v1.MembershipTierCreateRes, error) {
	item, err := service.Membership().AdminCreateTier(ctx, req)
	if err != nil {
		return nil, err
	}
	return &v1.MembershipTierCreateRes{MembershipTierItem: *item}, nil
}

func (c *ControllerV1) MembershipTierUpdate(ctx context.Context, req *v1.MembershipTierUpdateReq) (*v1.MembershipTierUpdateRes, error) {
	return &v1.MembershipTierUpdateRes{}, service.Membership().AdminUpdateTier(ctx, req)
}

func (c *ControllerV1) MembershipTierDelete(ctx context.Context, req *v1.MembershipTierDeleteReq) (*v1.MembershipTierDeleteRes, error) {
	return &v1.MembershipTierDeleteRes{}, service.Membership().AdminDeleteTier(ctx, req.ID)
}

func (c *ControllerV1) UserMembership(ctx context.Context, req *v1.UserMembershipReq) (*v1.UserMembershipRes, error) {
	userID, ok := middleware.GetCurrentUserID(ctx)
	if !ok || userID == 0 {
		return &v1.UserMembershipRes{Active: false}, nil
	}
	return service.Membership().GetUserMembership(ctx, userID)
}

func (c *ControllerV1) MembershipSubscriberList(ctx context.Context, req *v1.MembershipSubscriberListReq) (*v1.MembershipSubscriberListRes, error) {
	return service.Membership().AdminListSubscribers(ctx, req)
}
