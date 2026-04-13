package v1

import "github.com/gogf/gf/v2/frame/g"

// ── Order status constants ──────────────────────────────────────────────────

const (
	OrderStatusPending   = 1
	OrderStatusPaid      = 2
	OrderStatusCompleted = 3
	OrderStatusCancelled = 4
	OrderStatusRefunded  = 5
)

// ── Order item types ────────────────────────────────────────────────────────

const (
	ItemTypePostUnlock  = "post_unlock"
	ItemTypeDownload    = "download"
	ItemTypeMembership  = "membership"
	ItemTypeProduct     = "product"
	ItemTypeTopup       = "topup"
)

// ── Shared types ────────────────────────────────────────────────────────────

type OrderItem struct {
	ItemType  string `json:"item_type"`
	ItemID    string `json:"item_id"`
	Title     string `json:"title"`
	UnitPrice int    `json:"unit_price"` // cents
	Quantity  int    `json:"quantity"`
	Snapshot  string `json:"snapshot,omitempty"`
}

type OrderListItem struct {
	ID          int64       `json:"id"`
	OrderNo     string      `json:"order_no"`
	UserID      int64       `json:"user_id"`
	Status      int         `json:"status"`
	TotalAmount int         `json:"total_amount"`
	PaidAmount  int         `json:"paid_amount"`
	CreditsUsed int         `json:"credits_used"`
	BalanceUsed int         `json:"balance_used"`
	Currency    string      `json:"currency"`
	Items       []OrderItem `json:"items"`
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

// ── User API: Create Order ──────────────────────────────────────────────────

type OrderCreateReq struct {
	g.Meta     `path:"/api/v1/orders" method:"post" tags:"Order" summary:"Create order"`
	Items      []OrderItem `json:"items" v:"required"`
	UseCredits int         `json:"use_credits,omitempty"` // credits to apply
	UseBalance int         `json:"use_balance,omitempty"` // balance cents to apply
	Currency   string      `json:"currency,omitempty"`
}

type OrderCreateRes struct {
	ID      int64  `json:"id"`
	OrderNo string `json:"order_no"`
	Status  int    `json:"status"`
	Amount  int    `json:"amount"`
}

// ── User API: List Orders ───────────────────────────────────────────────────

type OrderListReq struct {
	g.Meta   `path:"/api/v1/orders" method:"get" tags:"Order" summary:"List orders"`
	Page     int  `json:"page" d:"1"`
	PageSize int  `json:"page_size" d:"20"`
	Status   *int `json:"status,omitempty"`
}

type OrderListRes struct {
	Data       []OrderListItem `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ── User API: Pay Order ─────────────────────────────────────────────────────

type OrderPayReq struct {
	g.Meta   `path:"/api/v1/orders/{id}/pay" method:"post" tags:"Order" summary:"Pay order"`
	ID       int64  `json:"-" v:"required"`
	Provider string `json:"provider" v:"required"` // "alipay","paypal","balance"
}

type OrderPayRes struct {
	PaymentURL string `json:"payment_url,omitempty"` // redirect URL for external providers
	QRCode     string `json:"qr_code,omitempty"`     // QR code content
	Paid       bool   `json:"paid"`                  // true if paid immediately (balance)
}

// ── User API: Cancel Order ──────────────────────────────────────────────────

type OrderCancelReq struct {
	g.Meta `path:"/api/v1/orders/{id}/cancel" method:"post" tags:"Order" summary:"Cancel order"`
	ID     int64 `json:"-" v:"required"`
}

type OrderCancelRes struct{}

// ── User API: Order Detail ──────────────────────────────────────────────────

type OrderDetailReq struct {
	g.Meta `path:"/api/v1/orders/{id}" method:"get" tags:"Order" summary:"Get order detail"`
	ID     int64 `json:"-" v:"required"`
}

type OrderDetailRes struct {
	OrderListItem
}

// ── User API: Check Purchase ────────────────────────────────────────────────

type PurchaseCheckReq struct {
	g.Meta     `path:"/api/v1/purchases/check" method:"get" tags:"Order" summary:"Check purchase status"`
	ObjectType string `json:"type" v:"required"`
	ObjectID   string `json:"id" v:"required"`
}

type PurchaseCheckRes struct {
	Purchased bool   `json:"purchased"`
	OrderID   *int64 `json:"order_id,omitempty"`
}

// ── Payment notify ──────────────────────────────────────────────────────────

type PaymentNotifyReq struct {
	g.Meta   `path:"/api/v1/payment/notify/{provider}" method:"post" tags:"Payment" summary:"Payment callback"`
	Provider string `json:"-" v:"required"`
}

type PaymentNotifyRes struct{}

// ── Admin API: Order Management ─────────────────────────────────────────────

type AdminOrderListReq struct {
	g.Meta   `path:"/admin/orders" method:"get" tags:"Admin Order" summary:"List all orders"`
	Page     int    `json:"page" d:"1"`
	PageSize int    `json:"page_size" d:"20"`
	Status   *int   `json:"status,omitempty"`
	UserID   *int64 `json:"user_id,omitempty"`
}

type AdminOrderListRes struct {
	Data       []OrderListItem `json:"data"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

// ── Admin API: Refund ───────────────────────────────────────────────────────

type OrderRefundReq struct {
	g.Meta `path:"/admin/orders/{id}/refund" method:"post" tags:"Admin Order" summary:"Refund order"`
	ID     int64  `json:"-" v:"required"`
	Reason string `json:"reason"`
}

type OrderRefundRes struct{}

// ── Admin API: Revenue Stats ────────────────────────────────────────────────

type RevenueStatsReq struct {
	g.Meta `path:"/admin/revenue/stats" method:"get" tags:"Admin Revenue" summary:"Revenue statistics"`
}

type RevenueStatsRes struct {
	TotalRevenue    int `json:"total_revenue"`    // cents
	TodayRevenue    int `json:"today_revenue"`    // cents
	TotalOrders     int `json:"total_orders"`
	TodayOrders     int `json:"today_orders"`
	ActiveMembers   int `json:"active_members"`
	PendingOrders   int `json:"pending_orders"`
}

// ── Membership Tier types ───────────────────────────────────────────────────

type MembershipTierItem struct {
	ID             int64  `json:"id"`
	Name           string `json:"name"`
	Slug           string `json:"slug"`
	Description    string `json:"description"`
	Price          int    `json:"price"` // cents
	DurationDays   int    `json:"duration_days"`
	DiscountPct    int    `json:"discount_pct"`
	AccessAll      bool   `json:"access_all"`
	CreditsMonthly int    `json:"credits_monthly"`
	Features       string `json:"features"` // JSON
	Status         int    `json:"status"`
	SortOrder      int    `json:"sort_order"`
}

type MembershipTierListReq struct {
	g.Meta `path:"/api/v1/membership/tiers" method:"get" tags:"Membership" summary:"List membership tiers"`
}

type MembershipTierListRes struct {
	Items []MembershipTierItem `json:"items"`
}

type MembershipTierCreateReq struct {
	g.Meta         `path:"/admin/membership/tiers" method:"post" tags:"Admin Membership" summary:"Create membership tier"`
	Name           string `json:"name" v:"required"`
	Slug           string `json:"slug" v:"required"`
	Description    string `json:"description"`
	Price          int    `json:"price" v:"required|min:0"`
	DurationDays   int    `json:"duration_days" v:"required|min:1"`
	DiscountPct    int    `json:"discount_pct" v:"between:0,100"`
	AccessAll      bool   `json:"access_all"`
	CreditsMonthly int    `json:"credits_monthly"`
	Features       string `json:"features"`
	Status         int    `json:"status" d:"1"`
	SortOrder      int    `json:"sort_order"`
}

type MembershipTierCreateRes struct {
	MembershipTierItem
}

type MembershipTierUpdateReq struct {
	g.Meta         `path:"/admin/membership/tiers/{id}" method:"put" tags:"Admin Membership" summary:"Update membership tier"`
	ID             int64  `json:"-" v:"required"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	Price          *int   `json:"price"`
	DurationDays   *int   `json:"duration_days"`
	DiscountPct    *int   `json:"discount_pct"`
	AccessAll      *bool  `json:"access_all"`
	CreditsMonthly *int   `json:"credits_monthly"`
	Features       string `json:"features"`
	Status         *int   `json:"status"`
	SortOrder      *int   `json:"sort_order"`
}

type MembershipTierUpdateRes struct{}

type MembershipTierDeleteReq struct {
	g.Meta `path:"/admin/membership/tiers/{id}" method:"delete" tags:"Admin Membership" summary:"Delete membership tier"`
	ID     int64 `json:"-" v:"required"`
}

type MembershipTierDeleteRes struct{}

// ── User Membership ─────────────────────────────────────────────────────────

type UserMembershipReq struct {
	g.Meta `path:"/api/v1/membership/me" method:"get" tags:"Membership" summary:"Get my membership"`
}

type UserMembershipRes struct {
	Active    bool                `json:"active"`
	Tier      *MembershipTierItem `json:"tier,omitempty"`
	ExpiresAt string              `json:"expires_at,omitempty"`
	AutoRenew bool                `json:"auto_renew"`
}

// ── Admin: Subscriber list ──────────────────────────────────────────────────

type MembershipSubscriberListReq struct {
	g.Meta   `path:"/admin/membership/subscribers" method:"get" tags:"Admin Membership" summary:"List subscribers"`
	Page     int `json:"page" d:"1"`
	PageSize int `json:"page_size" d:"20"`
}

type MembershipSubscriberListRes struct {
	Data       []SubscriberItem `json:"data"`
	Total      int              `json:"total"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	TotalPages int              `json:"total_pages"`
}

type SubscriberItem struct {
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	TierName  string `json:"tier_name"`
	Status    int    `json:"status"`
	StartedAt string `json:"started_at"`
	ExpiresAt string `json:"expires_at"`
	AutoRenew bool   `json:"auto_renew"`
}
