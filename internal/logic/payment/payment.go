package payment

import (
	"context"
	"encoding/json"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"

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

// ── Provider registry ─────────────────────────────────────────────────────────

var providers []PaymentProvider

func RegisterProvider(p PaymentProvider) {
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

// ── Service implementation ────────────────────────────────────────────────────

type sPayment struct{}

func New() service.IPayment { return &sPayment{} }

func init() {
	service.RegisterPayment(New())

	// Register built-in providers
	RegisterProvider(&alipayProvider{})
	RegisterProvider(&paypalProvider{})
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
