package payment

import v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"

type alipayProvider struct{}

func (p *alipayProvider) Slug() string  { return "alipay" }
func (p *alipayProvider) Label() string { return "支付宝" }
func (p *alipayProvider) Icon() string  { return "i-tabler-brand-alipay" }

func (p *alipayProvider) Fields() []v1.FieldDef {
	return []v1.FieldDef{
		{Key: "enabled", Label: "启用", Type: "switch"},
		{Key: "app_id", Label: "App ID", Type: "text", Required: true, Placeholder: "支付宝应用 APPID"},
		{Key: "private_key", Label: "应用私钥", Type: "password", Required: true, Placeholder: "RSA2 应用私钥"},
		{Key: "alipay_public_key", Label: "支付宝公钥", Type: "password", Required: true, Placeholder: "支付宝公钥"},
		{Key: "sign_type", Label: "签名类型", Type: "select", Options: []v1.FieldOption{
			{Label: "RSA2", Value: "RSA2"},
		}},
		{Key: "sandbox", Label: "沙箱模式", Type: "switch"},
	}
}

func (p *alipayProvider) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled":           false,
		"app_id":            "",
		"private_key":       "",
		"alipay_public_key": "",
		"sign_type":         "RSA2",
		"sandbox":           false,
	}
}

func (p *alipayProvider) MaskConfig(cfg map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(cfg))
	for k, v := range cfg {
		out[k] = v
	}
	out["private_key"] = maskKey(strVal(cfg, "private_key"))
	out["alipay_public_key"] = maskKey(strVal(cfg, "alipay_public_key"))
	return out
}

func (p *alipayProvider) MergeConfig(existing, incoming map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(incoming))
	for k, v := range incoming {
		merged[k] = v
	}
	// Preserve secrets when masked placeholder is submitted
	if isMasked(strVal(incoming, "private_key")) {
		merged["private_key"] = strVal(existing, "private_key")
	}
	if isMasked(strVal(incoming, "alipay_public_key")) {
		merged["alipay_public_key"] = strVal(existing, "alipay_public_key")
	}
	// Default sign_type
	if strVal(merged, "sign_type") == "" {
		merged["sign_type"] = "RSA2"
	}
	return merged
}
