package payment

import (
	"context"
	"encoding/json"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
	sdk "github.com/nuxtblog/nuxtblog/sdk"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ── Provider interface ────────────────────────────────────────────────────────

// PaymentProvider defines what each payment provider must implement.
type PaymentProvider interface {
	Slug() string
	Label() string
	Icon() string
	Fields() []v1.FieldDef
	DefaultConfig() map[string]interface{}
	// MaskConfig returns a copy with sensitive values masked for API response.
	MaskConfig(cfg map[string]interface{}) map[string]interface{}
	// MergeConfig merges incoming values into existing, preserving masked placeholders.
	MergeConfig(existing, incoming map[string]interface{}) map[string]interface{}
}

// PaymentGatewayInternal extends PaymentProvider with payment operations.
// Providers implementing this can create payments and verify callbacks.
type PaymentGatewayInternal interface {
	CreatePayment(ctx context.Context, cfg map[string]interface{}, req service.CreatePaymentReq) (*service.CreatePaymentRes, error)
	VerifyNotify(ctx context.Context, cfg map[string]interface{}, body []byte, headers map[string]string) (*service.NotifyResult, error)
}

// ── Provider registry ─────────────────────────────────────────────────────────

var providers []PaymentProvider

func RegisterProvider(p PaymentProvider) {
	// Replace if slug already exists (e.g. plugin overriding built-in)
	for i, existing := range providers {
		if existing.Slug() == p.Slug() {
			providers[i] = p
			return
		}
	}
	providers = append(providers, p)
}

func findProvider(slug string) PaymentProvider {
	for _, p := range providers {
		if p.Slug() == slug {
			return p
		}
	}
	return nil
}

// ── Plugin gateway adapter ────────────────────────────────────────────────────

// pluginPaymentAdapter adapts an sdk.PaymentGateway to internal PaymentProvider + PaymentGatewayInternal.
type pluginPaymentAdapter struct {
	gateway sdk.PaymentGateway
	info    sdk.PaymentProviderInfo
}

func (a *pluginPaymentAdapter) Slug() string  { return a.info.Slug }
func (a *pluginPaymentAdapter) Label() string { return a.info.Label }
func (a *pluginPaymentAdapter) Icon() string  { return a.info.Icon }

func (a *pluginPaymentAdapter) Fields() []v1.FieldDef {
	fields := make([]v1.FieldDef, len(a.info.Fields))
	for i, f := range a.info.Fields {
		opts := make([]v1.FieldOption, len(f.Options))
		for j, o := range f.Options {
			opts[j] = v1.FieldOption{Label: o.Label, Value: o.Value}
		}
		fields[i] = v1.FieldDef{
			Key:         f.Key,
			Label:       f.Label,
			Type:        f.Type,
			Required:    f.Required,
			Placeholder: f.Placeholder,
			Options:     opts,
		}
	}
	return fields
}

func (a *pluginPaymentAdapter) DefaultConfig() map[string]interface{} {
	return a.info.DefaultConfig
}

func (a *pluginPaymentAdapter) MaskConfig(cfg map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(cfg))
	for k, v := range cfg {
		out[k] = v
	}
	// Auto-mask fields with type="password"
	for _, f := range a.info.Fields {
		if f.Type == "password" {
			out[f.Key] = maskKey(strVal(cfg, f.Key))
		}
	}
	return out
}

func (a *pluginPaymentAdapter) MergeConfig(existing, incoming map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(incoming))
	for k, v := range incoming {
		merged[k] = v
	}
	// Preserve masked password fields
	for _, f := range a.info.Fields {
		if f.Type == "password" {
			if isMasked(strVal(incoming, f.Key)) {
				merged[f.Key] = strVal(existing, f.Key)
			}
		}
	}
	return merged
}

func (a *pluginPaymentAdapter) CreatePayment(ctx context.Context, cfg map[string]interface{}, req service.CreatePaymentReq) (*service.CreatePaymentRes, error) {
	sdkReq := sdk.PaymentRequest{
		OrderNo:   req.OrderNo,
		Amount:    req.Amount,
		Currency:  req.Currency,
		Subject:   req.Subject,
		NotifyURL: req.NotifyURL,
		ReturnURL: req.ReturnURL,
		ClientIP:  req.ClientIP,
	}
	res, err := a.gateway.CreatePayment(ctx, cfg, sdkReq)
	if err != nil {
		return nil, err
	}
	return &service.CreatePaymentRes{
		PaymentURL: res.PaymentURL,
		QRCode:     res.QRCode,
		Method:     res.Method,
	}, nil
}

func (a *pluginPaymentAdapter) VerifyNotify(ctx context.Context, cfg map[string]interface{}, body []byte, headers map[string]string) (*service.NotifyResult, error) {
	res, err := a.gateway.HandleNotify(ctx, cfg, body, headers)
	if err != nil {
		return nil, err
	}
	return &service.NotifyResult{
		Success:      res.Success,
		OrderNo:      res.OrderNo,
		Amount:       res.Amount,
		ProviderTxID: res.ProviderTxID,
		RawResponse:  res.RawResponse,
	}, nil
}

// RegisterPluginGateway registers an SDK PaymentGateway as an internal provider.
func RegisterPluginGateway(gw sdk.PaymentGateway) {
	adapter := &pluginPaymentAdapter{gateway: gw, info: gw.ProviderInfo()}
	RegisterProvider(adapter)
}

// ── Service implementation ────────────────────────────────────────────────────

type sPayment struct{}

func New() service.IPayment { return &sPayment{} }

func init() {
	service.RegisterPayment(New())

	// Only balance is built-in; alipay/wechat/paypal are provided by plugins
	RegisterProvider(&balanceProvider{})
}

func requireAdmin(ctx context.Context) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}
	return nil
}

// ── Options table helpers ─────────────────────────────────────────────────────

func optionKey(slug string) string {
	return "payment_" + slug
}

func loadConfig(ctx context.Context, slug string) (map[string]interface{}, error) {
	type row struct{ Value string `orm:"value"` }
	var r row
	_ = g.DB().Ctx(ctx).Model("options").Where("key", optionKey(slug)).Scan(&r)
	if r.Value == "" {
		return nil, nil
	}
	var cfg map[string]interface{}
	if err := json.Unmarshal([]byte(r.Value), &cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func saveConfig(ctx context.Context, slug string, cfg map[string]interface{}) error {
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	key := optionKey(slug)
	val := string(b)
	cnt, _ := g.DB().Ctx(ctx).Model("options").Where("key", key).Count()
	if cnt > 0 {
		_, err = g.DB().Ctx(ctx).Model("options").Where("key", key).
			Data(g.Map{"value": val}).Update()
	} else {
		_, err = g.DB().Ctx(ctx).Model("options").
			Data(g.Map{"key": key, "value": val, "autoload": 0}).Insert()
	}
	return err
}

// ── Build ProviderInfo ────────────────────────────────────────────────────────

func buildProviderInfo(p PaymentProvider, cfg map[string]interface{}) v1.ProviderInfo {
	if cfg == nil {
		cfg = p.DefaultConfig()
	}
	enabled, _ := cfg["enabled"].(bool)
	masked := p.MaskConfig(cfg)
	return v1.ProviderInfo{
		Slug:    p.Slug(),
		Label:   p.Label(),
		Icon:    p.Icon(),
		Enabled: enabled,
		Config:  masked,
		Fields:  p.Fields(),
	}
}

// ── IPayment methods ──────────────────────────────────────────────────────────

func (s *sPayment) ListProviders(ctx context.Context) (*v1.PaymentListProvidersRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	items := make([]v1.ProviderInfo, 0, len(providers))
	for _, p := range providers {
		cfg, _ := loadConfig(ctx, p.Slug())
		items = append(items, buildProviderInfo(p, cfg))
	}
	return &v1.PaymentListProvidersRes{Items: items}, nil
}

func (s *sPayment) GetProviderConfig(ctx context.Context, slug string) (*v1.PaymentGetProviderConfigRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	p := findProvider(slug)
	if p == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment provider not found: "+slug)
	}
	cfg, err := loadConfig(ctx, slug)
	if err != nil {
		return nil, err
	}
	info := buildProviderInfo(p, cfg)
	return &v1.PaymentGetProviderConfigRes{ProviderInfo: info}, nil
}

func (s *sPayment) SetProviderConfig(ctx context.Context, slug string, incoming map[string]interface{}) (*v1.PaymentSetProviderConfigRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	p := findProvider(slug)
	if p == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment provider not found: "+slug)
	}
	existing, err := loadConfig(ctx, slug)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		existing = p.DefaultConfig()
	}

	merged := p.MergeConfig(existing, incoming)
	if err := saveConfig(ctx, slug, merged); err != nil {
		return nil, err
	}

	info := buildProviderInfo(p, merged)
	return &v1.PaymentSetProviderConfigRes{ProviderInfo: info}, nil
}

func (s *sPayment) ListEnabledProviders(ctx context.Context) ([]v1.ProviderBasicInfo, error) {
	var result []v1.ProviderBasicInfo
	for _, p := range providers {
		cfg, _ := loadConfig(ctx, p.Slug())
		if cfg == nil {
			cfg = p.DefaultConfig()
		}
		enabled, _ := cfg["enabled"].(bool)
		if enabled {
			result = append(result, v1.ProviderBasicInfo{
				Slug:  p.Slug(),
				Label: p.Label(),
				Icon:  p.Icon(),
			})
		}
	}
	return result, nil
}

func (s *sPayment) CreatePayment(ctx context.Context, provider string, req service.CreatePaymentReq) (*service.CreatePaymentRes, error) {
	p := findProvider(provider)
	if p == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment provider not found: "+provider)
	}
	gw, ok := p.(PaymentGatewayInternal)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeNotSupported, "provider does not support payment creation: "+provider)
	}
	cfg, err := loadConfig(ctx, provider)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = p.DefaultConfig()
	}
	return gw.CreatePayment(ctx, cfg, req)
}

func (s *sPayment) VerifyNotify(ctx context.Context, provider string, body []byte, headers map[string]string) (*service.NotifyResult, error) {
	p := findProvider(provider)
	if p == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "payment provider not found: "+provider)
	}
	gw, ok := p.(PaymentGatewayInternal)
	if !ok {
		return nil, gerror.NewCode(gcode.CodeNotSupported, "provider does not support notify verification: "+provider)
	}
	cfg, err := loadConfig(ctx, provider)
	if err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = p.DefaultConfig()
	}
	return gw.VerifyNotify(ctx, cfg, body, headers)
}

// ── Shared mask helper ────────────────────────────────────────────────────────

func maskKey(key string) string {
	if len(key) > 8 {
		return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
	}
	if key != "" {
		return "****"
	}
	return ""
}

func isMasked(val string) bool {
	return strings.Contains(val, "****")
}

func strVal(m map[string]interface{}, key string) string {
	v, _ := m[key].(string)
	return v
}

func boolVal(m map[string]interface{}, key string) bool {
	v, _ := m[key].(bool)
	return v
}
