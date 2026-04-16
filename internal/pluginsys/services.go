package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	sdk "github.com/nuxtblog/nuxtblog/sdk"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

// ─── Platform service implementations ───────────────────────────────────────

// pluginDB implements sdk.DB
type pluginDB struct {
	pluginID string
	prefix   string
	caps     *DBCap
	trust    TrustLevel
}

func newPluginDB(pluginID string, caps *DBCap, trust TrustLevel) *pluginDB {
	return &pluginDB{
		pluginID: pluginID,
		prefix:   sanitizeTablePrefix(pluginID),
		caps:     caps,
		trust:    trust,
	}
}

func (d *pluginDB) Query(sql string, args ...any) ([]map[string]any, error) {
	if err := (&sqlGuard{d.prefix, d.caps, d.trust}).validate(sql); err != nil {
		return nil, fmt.Errorf("db access denied: %w", err)
	}
	ctx := context.Background()
	result, err := g.DB().GetAll(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(result))
	for _, r := range result {
		rows = append(rows, r.Map())
	}
	return rows, nil
}

func (d *pluginDB) Execute(sql string, args ...any) (int64, error) {
	if err := (&sqlGuard{d.prefix, d.caps, d.trust}).validate(sql); err != nil {
		return 0, fmt.Errorf("db access denied: %w", err)
	}
	ctx := context.Background()
	res, err := g.DB().Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// pluginStore implements sdk.Store
type pluginStore struct {
	pluginID string
}

func newPluginStore(pluginID string) *pluginStore {
	return &pluginStore{pluginID: pluginID}
}

func (s *pluginStore) storeKey(key string) string {
	return "plugin_store:" + s.pluginID + ":" + key
}

func (s *pluginStore) Get(key string) (any, error) {
	val, err := g.DB().Ctx(context.Background()).
		Model("options").Where("key", s.storeKey(key)).Value("value")
	if err != nil {
		return nil, err
	}
	if val.IsNil() {
		return nil, nil
	}
	var v any
	if json.Unmarshal([]byte(val.String()), &v) == nil {
		return v, nil
	}
	return val.String(), nil
}

func (s *pluginStore) Set(key string, value any) error {
	raw, _ := json.Marshal(value)
	return upsertOption(context.Background(), s.storeKey(key), string(raw))
}

func (s *pluginStore) Delete(key string) error {
	_, err := g.DB().Ctx(context.Background()).Model("options").
		Where("key", s.storeKey(key)).Delete()
	return err
}

func (s *pluginStore) Increment(key string, delta ...int64) (int64, error) {
	d := int64(1)
	if len(delta) > 0 {
		d = delta[0]
	}
	ctx := context.Background()
	db := g.DB()
	dbKey := s.storeKey(key)

	dbType := ""
	if cfg := db.GetConfig(); cfg != nil {
		dbType = strings.ToLower(cfg.Type)
	}
	var sql string
	switch dbType {
	case "mysql", "mariadb":
		sql = "INSERT INTO `options` (`key`, `value`, `autoload`) VALUES (?, ?, 0)" +
			" ON DUPLICATE KEY UPDATE `value` = CAST(`value` AS SIGNED) + CAST(VALUES(`value`) AS SIGNED)"
	default:
		sql = "INSERT INTO options (key, value, autoload) VALUES (?, ?, 0)" +
			" ON CONFLICT(key) DO UPDATE SET value = CAST(options.value AS INTEGER) + CAST(excluded.value AS INTEGER)"
	}
	raw, _ := json.Marshal(d)
	if _, err := db.Exec(ctx, sql, dbKey, string(raw)); err != nil {
		return 0, err
	}
	val, _ := db.Ctx(ctx).Model("options").Where("key", dbKey).Value("value")
	var result int64
	if val != nil && !val.IsNil() {
		_ = json.Unmarshal([]byte(val.String()), &result)
	}
	return result, nil
}

func (s *pluginStore) DeletePrefix(prefix string) (int64, error) {
	dbPrefix := "plugin_store:" + s.pluginID + ":" + prefix
	res, err := g.DB().Ctx(context.Background()).Model("options").
		WhereLike("key", dbPrefix+"%").Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// upsertOption writes key=value into the options table atomically.
func upsertOption(ctx context.Context, key, value string) error {
	db := g.DB()
	dbType := ""
	if cfg := db.GetConfig(); cfg != nil {
		dbType = strings.ToLower(cfg.Type)
	}
	var sql string
	switch dbType {
	case "mysql", "mariadb":
		sql = "INSERT INTO `options` (`key`, `value`, `autoload`) VALUES (?, ?, 0)" +
			" ON DUPLICATE KEY UPDATE `value` = VALUES(`value`)"
	default:
		sql = "INSERT INTO options (key, value, autoload) VALUES (?, ?, 0)" +
			" ON CONFLICT(key) DO UPDATE SET value = excluded.value"
	}
	_, err := db.Exec(ctx, sql, key, value)
	return err
}

// pluginSettings implements sdk.Settings with a 30-second TTL cache
// to avoid hitting the DB on every call (e.g. inside high-frequency filter handlers).
type pluginSettings struct {
	pluginID string

	mu      sync.Mutex
	cached  map[string]any
	cachedAt time.Time
}

const settingsCacheTTL = 30 * time.Second

func newPluginSettings(pluginID string) *pluginSettings {
	return &pluginSettings{pluginID: pluginID}
}

// InvalidateCache forces the next GetAll to re-read from the database.
// Called by ReactivatePlugin after settings may have changed.
func (s *pluginSettings) InvalidateCache() {
	s.mu.Lock()
	s.cached = nil
	s.cachedAt = time.Time{}
	s.mu.Unlock()
}

func (s *pluginSettings) Get(key string) any {
	return s.GetAll()[key]
}

func (s *pluginSettings) GetAll() map[string]any {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.cached != nil && time.Since(s.cachedAt) < settingsCacheTTL {
		return s.cached
	}

	val, _ := g.DB().Ctx(context.Background()).
		Model("plugins").Where("id", s.pluginID).Value("settings")
	var m map[string]any
	if val != nil && !val.IsNil() {
		_ = json.Unmarshal([]byte(val.String()), &m)
	}
	if m == nil {
		m = map[string]any{}
	}
	s.cached = m
	s.cachedAt = time.Now()
	return m
}

// ─── Commerce service ─────────────────────────────────────────────────────────

// pluginCommerce implements sdk.Commerce by proxying to registered provider plugins
// for wallet/credits/membership/entitlement, and to the core payment service for payments.
type pluginCommerce struct{ pluginID string }

func newPluginCommerce(pluginID string) *pluginCommerce {
	return &pluginCommerce{pluginID: pluginID}
}

// ── Payment gateway (core infrastructure) ──

func (c *pluginCommerce) GetEnabledPaymentMethods(ctx context.Context) ([]sdk.PaymentMethod, error) {
	svc := service.Payment()
	if svc == nil {
		return nil, fmt.Errorf("payment service not registered")
	}
	providers, err := svc.ListEnabledProviders(ctx)
	if err != nil {
		return nil, err
	}
	var methods []sdk.PaymentMethod
	for _, p := range providers {
		methods = append(methods, sdk.PaymentMethod{
			Slug:  p.Slug,
			Label: p.Label,
			Icon:  p.Icon,
		})
	}
	return methods, nil
}

func (c *pluginCommerce) CreatePayment(ctx context.Context, provider string, req sdk.PaymentRequest) (*sdk.PaymentResult, error) {
	svc := service.Payment()
	if svc == nil {
		return nil, fmt.Errorf("payment service not registered")
	}
	res, err := svc.CreatePayment(ctx, provider, service.CreatePaymentReq{
		OrderNo:   req.OrderNo,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Subject:   req.Subject,
		NotifyURL: req.NotifyURL,
		ReturnURL: req.ReturnURL,
		ClientIP:  req.ClientIP,
	})
	if err != nil {
		return nil, err
	}
	return &sdk.PaymentResult{
		PaymentURL: res.PaymentURL,
		QRCode:     res.QRCode,
		Method:     res.Method,
	}, nil
}

func (c *pluginCommerce) VerifyNotify(ctx context.Context, provider string, body []byte, headers map[string]string) (*sdk.PaymentNotifyResult, error) {
	svc := service.Payment()
	if svc == nil {
		return nil, fmt.Errorf("payment service not registered")
	}
	res, err := svc.VerifyNotify(ctx, provider, body, headers)
	if err != nil {
		return nil, err
	}
	return &sdk.PaymentNotifyResult{
		Success:      res.Success,
		OrderNo:      res.OrderNo,
		Amount:       res.Amount,
		ProviderTxID: res.ProviderTxID,
		RawResponse:  res.RawResponse,
	}, nil
}

// ── Wallet (proxied to WalletProvider plugin) ──

func (c *pluginCommerce) GetUserBalance(ctx context.Context, userID int64) (int, int, error) {
	var balance, credits int
	if wp := defaultMgr.WalletProvider(); wp != nil {
		b, _ := wp.GetBalance(ctx, userID)
		balance = b
	}
	if cp := defaultMgr.CreditsProvider(); cp != nil {
		cr, _ := cp.GetBalance(ctx, userID)
		credits = cr
	}
	return balance, credits, nil
}

func (c *pluginCommerce) SpendBalance(ctx context.Context, userID int64, amount int, refType, refID, note string) (int, error) {
	wp := defaultMgr.WalletProvider()
	if wp == nil {
		return 0, fmt.Errorf("no wallet provider registered")
	}
	return wp.Spend(ctx, userID, amount, refType, refID, note)
}

func (c *pluginCommerce) RefundBalance(ctx context.Context, userID int64, amount int, refType, refID, note string) (int, error) {
	wp := defaultMgr.WalletProvider()
	if wp == nil {
		return 0, fmt.Errorf("no wallet provider registered")
	}
	return wp.Refund(ctx, userID, amount, refType, refID, note)
}

func (c *pluginCommerce) TopupWallet(ctx context.Context, userID int64, amount int, provider, txID string) error {
	wp := defaultMgr.WalletProvider()
	if wp == nil {
		return fmt.Errorf("no wallet provider registered")
	}
	return wp.Topup(ctx, userID, amount, provider, txID)
}

// ── Credits (proxied to CreditsProvider plugin) ──

func (c *pluginCommerce) SpendCredits(ctx context.Context, userID int64, amount int, refType, refID, note string) (int, error) {
	cp := defaultMgr.CreditsProvider()
	if cp == nil {
		return 0, fmt.Errorf("no credits provider registered")
	}
	return cp.Spend(ctx, userID, amount, refType, refID, note)
}

func (c *pluginCommerce) EarnCredits(ctx context.Context, userID int64, amount int, source, refID, note string) (int, error) {
	cp := defaultMgr.CreditsProvider()
	if cp == nil {
		return 0, fmt.Errorf("no credits provider registered")
	}
	return cp.Earn(ctx, userID, amount, source, refID, note)
}

// ── Entitlement (proxied to EntitlementProvider plugin) ──

func (c *pluginCommerce) GrantEntitlement(ctx context.Context, userID int64, objectType, objectID string) error {
	ep := defaultMgr.EntitlementProvider()
	if ep == nil {
		return fmt.Errorf("no entitlement provider registered")
	}
	return ep.Grant(ctx, userID, objectType, objectID)
}

func (c *pluginCommerce) RevokeEntitlement(ctx context.Context, userID int64, objectType, objectID string) error {
	ep := defaultMgr.EntitlementProvider()
	if ep == nil {
		return fmt.Errorf("no entitlement provider registered")
	}
	return ep.Revoke(ctx, userID, objectType, objectID)
}

func (c *pluginCommerce) CheckEntitlement(ctx context.Context, userID int64, objectType, objectID string) (bool, error) {
	ep := defaultMgr.EntitlementProvider()
	if ep == nil {
		return false, nil // no provider = no entitlement
	}
	return ep.Check(ctx, userID, objectType, objectID)
}

// ── Membership (proxied to MembershipProvider plugin) ──

func (c *pluginCommerce) ActivateMembership(ctx context.Context, userID int64, tierID int64) error {
	mp := defaultMgr.MembershipProvider()
	if mp == nil {
		return fmt.Errorf("no membership provider registered")
	}
	return mp.Activate(ctx, userID, tierID)
}

func (c *pluginCommerce) CheckMembershipAccess(ctx context.Context, userID int64) (bool, error) {
	mp := defaultMgr.MembershipProvider()
	if mp == nil {
		return false, nil
	}
	return mp.CheckAccess(ctx, userID)
}

func (c *pluginCommerce) ListMembershipTiers(ctx context.Context) ([]sdk.MembershipTier, error) {
	mp := defaultMgr.MembershipProvider()
	if mp == nil {
		return nil, nil
	}
	return mp.ListTiers(ctx)
}

func (c *pluginCommerce) GetUserMembershipTier(ctx context.Context, userID int64) (*sdk.MembershipTier, error) {
	mp := defaultMgr.MembershipProvider()
	if mp == nil {
		return nil, nil
	}
	return mp.GetUserTier(ctx, userID)
}

func (c *pluginCommerce) GetCreditsExchangeRate(ctx context.Context) (int, error) {
	cp := defaultMgr.CreditsProvider()
	if cp == nil {
		return 0, nil
	}
	return cp.ExchangeRate(ctx)
}

// pluginLogger implements sdk.Logger
type pluginLogger struct {
	pluginID string
}

func newPluginLogger(pluginID string) *pluginLogger {
	return &pluginLogger{pluginID: pluginID}
}

func (l *pluginLogger) Info(msg string) {
	g.Log().Infof(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}

func (l *pluginLogger) Warn(msg string) {
	g.Log().Warningf(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}

func (l *pluginLogger) Error(msg string) {
	g.Log().Errorf(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}

func (l *pluginLogger) Debug(msg string) {
	g.Log().Debugf(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}
