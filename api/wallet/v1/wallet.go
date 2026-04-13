package v1

import "github.com/gogf/gf/v2/frame/g"

// ── Wallet balance ──────────────────────────────────────────────────────────

type WalletBalanceReq struct {
	g.Meta `path:"/api/v1/wallet/balance" method:"get" tags:"Wallet" summary:"Get wallet balance"`
}

type WalletBalanceRes struct {
	Balance    int `json:"balance"`     // cents
	Frozen     int `json:"frozen"`      // cents
	TotalTopup int `json:"total_topup"` // cents
	TotalSpent int `json:"total_spent"` // cents
	Credits    int `json:"credits"`
}

// ── Wallet ledger ───────────────────────────────────────────────────────────

type WalletLedgerReq struct {
	g.Meta   `path:"/api/v1/wallet/ledger" method:"get" tags:"Wallet" summary:"Wallet ledger"`
	Page     int `json:"page" d:"1"`
	PageSize int `json:"page_size" d:"20"`
}

type WalletLedgerRes struct {
	Data       []LedgerItem `json:"data"`
	Total      int          `json:"total"`
	Page       int          `json:"page"`
	PageSize   int          `json:"page_size"`
	TotalPages int          `json:"total_pages"`
}

type LedgerItem struct {
	ID            int64  `json:"id"`
	Type          int    `json:"type"`
	Amount        int    `json:"amount"`
	BalanceAfter  int    `json:"balance_after"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   string `json:"reference_id"`
	Note          string `json:"note"`
	CreatedAt     string `json:"created_at"`
}

// ── Wallet topup ────────────────────────────────────────────────────────────

type WalletTopupReq struct {
	g.Meta   `path:"/api/v1/wallet/topup" method:"post" tags:"Wallet" summary:"Topup wallet"`
	Amount   int    `json:"amount" v:"required|min:1"` // cents
	Provider string `json:"provider"`                  // payment provider slug
}

type WalletTopupRes struct {
	OrderID    int64  `json:"order_id"`
	OrderNo    string `json:"order_no"`
	PaymentURL string `json:"payment_url,omitempty"`
}

// ── Admin: Wallet adjust ────────────────────────────────────────────────────

type WalletAdminAdjustReq struct {
	g.Meta `path:"/admin/wallet/adjust" method:"post" tags:"Admin Wallet" summary:"Adjust user wallet"`
	UserID int64  `json:"user_id" v:"required"`
	Amount int    `json:"amount" v:"required"` // signed: positive=add, negative=deduct
	Note   string `json:"note"`
}

type WalletAdminAdjustRes struct{}

// ── Credits ledger ──────────────────────────────────────────────────────────

type CreditsLedgerReq struct {
	g.Meta   `path:"/api/v1/credits/ledger" method:"get" tags:"Credits" summary:"Credits ledger"`
	Page     int `json:"page" d:"1"`
	PageSize int `json:"page_size" d:"20"`
}

type CreditsLedgerRes struct {
	Data       []CreditsLedgerItem `json:"data"`
	Total      int                 `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}

type CreditsLedgerItem struct {
	ID            int64  `json:"id"`
	Type          int    `json:"type"`
	Amount        int    `json:"amount"`
	BalanceAfter  int    `json:"balance_after"`
	Source        string `json:"source"`
	ReferenceType string `json:"reference_type"`
	ReferenceID   string `json:"reference_id"`
	Note          string `json:"note"`
	CreatedAt     string `json:"created_at"`
}

// ── Admin: Credits adjust ───────────────────────────────────────────────────

type CreditsAdminAdjustReq struct {
	g.Meta `path:"/admin/credits/adjust" method:"post" tags:"Admin Credits" summary:"Adjust user credits"`
	UserID int64  `json:"user_id" v:"required"`
	Amount int    `json:"amount" v:"required"` // signed
	Note   string `json:"note"`
}

type CreditsAdminAdjustRes struct{}
