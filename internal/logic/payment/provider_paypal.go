package payment

import v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"

type paypalProvider struct{}

func (p *paypalProvider) Slug() string  { return "paypal" }
func (p *paypalProvider) Label() string { return "PayPal" }
func (p *paypalProvider) Icon() string  { return "i-tabler-brand-paypal" }

func (p *paypalProvider) Fields() []v1.FieldDef {
	return []v1.FieldDef{
		{Key: "enabled", Label: "启用", Type: "switch"},
		{Key: "client_id", Label: "Client ID", Type: "text", Required: true, Placeholder: "PayPal Client ID"},
		{Key: "client_secret", Label: "Client Secret", Type: "password", Required: true, Placeholder: "PayPal Client Secret"},
		{Key: "mode", Label: "环境", Type: "select", Options: []v1.FieldOption{
			{Label: "Sandbox", Value: "sandbox"},
			{Label: "Live", Value: "live"},
		}},
		{Key: "currency", Label: "币种", Type: "select", Options: []v1.FieldOption{
			{Label: "USD", Value: "USD"},
			{Label: "EUR", Value: "EUR"},
			{Label: "GBP", Value: "GBP"},
			{Label: "JPY", Value: "JPY"},
			{Label: "CNY", Value: "CNY"},
		}},
	}
}

func (p *paypalProvider) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled":       false,
		"client_id":     "",
		"client_secret": "",
		"mode":          "sandbox",
		"currency":      "USD",
	}
}

func (p *paypalProvider) MaskConfig(cfg map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(cfg))
	for k, v := range cfg {
		out[k] = v
	}
	out["client_secret"] = maskKey(strVal(cfg, "client_secret"))
	return out
}

func (p *paypalProvider) MergeConfig(existing, incoming map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(incoming))
	for k, v := range incoming {
		merged[k] = v
	}
	if isMasked(strVal(incoming, "client_secret")) {
		merged["client_secret"] = strVal(existing, "client_secret")
	}
	if strVal(merged, "mode") == "" {
		merged["mode"] = "sandbox"
	}
	if strVal(merged, "currency") == "" {
		merged["currency"] = "USD"
	}
	return merged
}
