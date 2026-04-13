package payload

// OrderCreated is delivered when a new order is placed.
type OrderCreated struct {
	OrderID  int64
	OrderNo  string
	UserID   int64
	Amount   int // cents
	ItemType string
	ItemID   string
}

// OrderPaid is delivered when payment is confirmed for an order.
type OrderPaid struct {
	OrderID    int64
	OrderNo    string
	UserID     int64
	Amount     int
	Provider   string
	ProviderTx string
}

// OrderCompleted is delivered when an order reaches final completion.
type OrderCompleted struct {
	OrderID  int64
	OrderNo  string
	UserID   int64
	ItemType string
	ItemID   string
}

// OrderCancelled is delivered when an order is cancelled.
type OrderCancelled struct {
	OrderID int64
	OrderNo string
	UserID  int64
}

// OrderRefunded is delivered when an order is refunded.
type OrderRefunded struct {
	OrderID int64
	OrderNo string
	UserID  int64
	Amount  int
}

// WalletTopup is delivered when a user's wallet is topped up.
type WalletTopup struct {
	UserID       int64
	Amount       int // cents, positive
	BalanceAfter int
	OrderID      int64
}

// WalletSpend is delivered when balance is deducted from a user's wallet.
type WalletSpend struct {
	UserID       int64
	Amount       int // cents, positive
	BalanceAfter int
	OrderID      int64
}

// CreditsEarned is delivered when a user earns credits.
type CreditsEarned struct {
	UserID       int64
	Amount       int
	BalanceAfter int
	Source       string // "checkin", "comment", "purchase", "admin"
}

// CreditsSpent is delivered when a user spends credits.
type CreditsSpent struct {
	UserID       int64
	Amount       int
	BalanceAfter int
	OrderID      int64
}

// MembershipActivated is delivered when a user's membership becomes active.
type MembershipActivated struct {
	UserID   int64
	TierID   int64
	TierSlug string
	OrderID  int64
}

// MembershipExpired is delivered when a user's membership expires.
type MembershipExpired struct {
	UserID   int64
	TierID   int64
	TierSlug string
}

// MembershipRenewed is delivered when a membership is renewed.
type MembershipRenewed struct {
	UserID   int64
	TierID   int64
	TierSlug string
	OrderID  int64
}

// ContentUnlocked is delivered when a user unlocks paid content.
type ContentUnlocked struct {
	UserID     int64
	ObjectType string // "post", "download"
	ObjectID   string
	OrderID    int64
}
