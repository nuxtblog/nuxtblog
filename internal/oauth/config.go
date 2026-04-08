package oauth

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
)

// ProviderConfig holds credentials and settings for one OAuth provider.
// Persisted in the options table as JSON under key "oauth_{name}".
type ProviderConfig struct {
	Enabled      bool   `json:"enabled"`
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	CallbackUrl  string `json:"callbackUrl"`
}

// GetConfig reads provider config from the options table (key: "oauth_{name}").
// Falls back to config.yaml values for backward compatibility.
func GetConfig(ctx context.Context, name string) ProviderConfig {
	type optRow struct {
		Value string `orm:"value"`
	}
	var row optRow
	_ = g.DB().Ctx(ctx).Model("options").
		Where("key", "oauth_"+name).
		Fields("value").
		Scan(&row)

	if row.Value != "" && row.Value != "null" {
		var cfg ProviderConfig
		if err := json.Unmarshal([]byte(row.Value), &cfg); err == nil {
			return cfg
		}
	}

	// Fallback: read from config.yaml (backward compatibility)
	prefix := fmt.Sprintf("auth.oauth.%s.", name)
	enabled, _ := g.Cfg().Get(ctx, prefix+"enabled")
	clientId, _ := g.Cfg().Get(ctx, prefix+"clientId")
	clientSecret, _ := g.Cfg().Get(ctx, prefix+"clientSecret")
	callbackUrl, _ := g.Cfg().Get(ctx, prefix+"callbackUrl")
	return ProviderConfig{
		Enabled:      enabled.Bool(),
		ClientId:     clientId.String(),
		ClientSecret: clientSecret.String(),
		CallbackUrl:  callbackUrl.String(),
	}
}
