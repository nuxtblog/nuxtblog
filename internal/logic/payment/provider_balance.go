package payment

import v1 "github.com/nuxtblog/nuxtblog/api/payment/v1"

// balanceProvider represents the wallet balance payment method.
// This provider is always enabled when the user has sufficient balance.
type balanceProvider struct{}

func (p *balanceProvider) Slug() string  { return "balance" }
func (p *balanceProvider) Label() string { return "余额支付" }
func (p *balanceProvider) Icon() string  { return "i-tabler-wallet" }

func (p *balanceProvider) Fields() []v1.FieldDef {
	return []v1.FieldDef{
		{Key: "enabled", Label: "启用", Type: "switch"},
	}
}

func (p *balanceProvider) DefaultConfig() map[string]interface{} {
	return map[string]interface{}{
		"enabled": true,
	}
}

func (p *balanceProvider) MaskConfig(cfg map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(cfg))
	for k, v := range cfg {
		out[k] = v
	}
	return out
}

func (p *balanceProvider) MergeConfig(existing, incoming map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{}, len(incoming))
	for k, v := range incoming {
		merged[k] = v
	}
	return merged
}
