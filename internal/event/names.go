package event

// Standard event names used throughout the application.
// Format: "<resource>.<action>"
//
// Consumers should use these constants rather than raw strings so that
// typos are caught at compile time.
const (
	// Post lifecycle
	PostCreated   = "post.created"
	PostUpdated   = "post.updated"
	PostPublished = "post.published"
	PostDeleted   = "post.deleted"
	PostViewed    = "post.viewed"

	// Comment lifecycle
	CommentCreated       = "comment.created"
	CommentDeleted       = "comment.deleted"
	CommentStatusChanged = "comment.status_changed"
	CommentApproved      = "comment.approved"

	// User lifecycle
	UserRegistered = "user.registered"
	UserUpdated    = "user.updated"
	UserDeleted    = "user.deleted"
	UserLoggedIn   = "user.login"
	UserLoggedOut  = "user.logout"

	// Social
	UserFollowed = "user.followed"

	// Engagement
	ReactionAdded   = "reaction.added"
	ReactionRemoved = "reaction.removed"

	// Media lifecycle
	MediaUploaded = "media.uploaded"
	MediaDeleted  = "media.deleted"

	// Taxonomy / term lifecycle
	TaxonomyCreated = "taxonomy.created"
	TaxonomyDeleted = "taxonomy.deleted"
	TermCreated     = "term.created"
	TermDeleted     = "term.deleted"

	// Checkin
	CheckinDone = "checkin.done"

	// Site settings
	OptionUpdated = "option.updated"

	// Plugin lifecycle
	PluginInstalled   = "plugin.installed"
	PluginUninstalled = "plugin.uninstalled"

	// Commerce: Orders
	OrderCreated   = "order.created"
	OrderPaid      = "order.paid"
	OrderCompleted = "order.completed"
	OrderCancelled = "order.cancelled"
	OrderRefunded  = "order.refunded"

	// Commerce: Wallet
	WalletTopup = "wallet.topup"
	WalletSpend = "wallet.spend"

	// Commerce: Credits
	CreditsEarned = "credits.earned"
	CreditsSpent  = "credits.spent"

	// Commerce: Membership
	MembershipActivated = "membership.activated"
	MembershipExpired   = "membership.expired"
	MembershipRenewed   = "membership.renewed"

	// Commerce: Content
	ContentUnlocked = "content.unlocked"
)
