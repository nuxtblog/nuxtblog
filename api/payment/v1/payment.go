package v1

import "github.com/gogf/gf/v2/frame/g"

// FieldDef describes a single config field so the frontend can render forms dynamically.
type FieldDef struct {
	Key         string            `json:"key"`
	Label       string            `json:"label"`
	Type        string            `json:"type"`                  // "text" | "password" | "switch" | "select"
	Required    bool              `json:"required,omitempty"`
	Placeholder string            `json:"placeholder,omitempty"`
	Options     []FieldOption     `json:"options,omitempty"`     // for "select" type
}

type FieldOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// ProviderInfo is the shape returned to the frontend for each registered payment provider.
type ProviderInfo struct {
	Slug    string                 `json:"slug"`
	Label   string                 `json:"label"`
	Icon    string                 `json:"icon"`
	Enabled bool                   `json:"enabled"`
	Config  map[string]interface{} `json:"config"`
	Fields  []FieldDef             `json:"fields"`
}

// ── API endpoints ─────────────────────────────────────────────────────────────

type PaymentListProvidersReq struct {
	g.Meta `path:"/admin/payment/providers" method:"get" tags:"Payment" summary:"List payment providers" auth:"true"`
}
type PaymentListProvidersRes struct {
	Items []ProviderInfo `json:"items"`
}

type PaymentGetProviderConfigReq struct {
	g.Meta `path:"/admin/payment/providers/{slug}/config" method:"get" tags:"Payment" summary:"Get provider config" auth:"true"`
	Slug   string `json:"-" v:"required"`
}
type PaymentGetProviderConfigRes struct {
	ProviderInfo
}

type PaymentSetProviderConfigReq struct {
	g.Meta `path:"/admin/payment/providers/{slug}/config" method:"put" tags:"Payment" summary:"Set provider config" auth:"true"`
	Slug   string                 `json:"-" v:"required"`
	Config map[string]interface{} `json:"config"`
}
type PaymentSetProviderConfigRes struct {
	ProviderInfo
}
