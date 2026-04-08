package notify

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nuxtblog/nuxtblog/internal/consts"

	"github.com/nuxtblog/nuxtblog/internal/notify/smsprovider"
)

func init() {
	Register(&SMSChannel{})
}

// SMSConfig is stored in the options table under key "notify_sms".
type SMSConfig struct {
	Provider        string `json:"provider"`          // "aliyun" | "twilio" | "tencent" | ""
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	SignName        string `json:"sign_name"`
	TemplateCode    string `json:"template_code"`
}

type SMSChannel struct{}

func (c *SMSChannel) Name() string { return consts.NotifyChannelSMS }

func (c *SMSChannel) Enabled(ctx context.Context) bool {
	cfg := loadSMSConfig(ctx)
	return cfg != nil && cfg.Provider != ""
}

func (c *SMSChannel) Send(ctx context.Context, msg Message) error {
	cfg := loadSMSConfig(ctx)
	if cfg == nil || cfg.Provider == "" {
		return nil
	}
	p, ok := smsprovider.Get(cfg.Provider)
	if !ok {
		g.Log().Warningf(ctx, "[notify] unknown SMS provider: %s", cfg.Provider)
		return nil
	}
	provCfg := smsprovider.Config{
		Provider:        cfg.Provider,
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		SignName:        cfg.SignName,
		TemplateCode:    cfg.TemplateCode,
	}
	return p.Send(ctx, provCfg, "", msg.Content)
}

func loadSMSConfig(ctx context.Context) *SMSConfig {
	val, err := g.DB().Ctx(ctx).Model("options").Where("key", "notify_sms").Value("value")
	if err != nil || val.IsEmpty() {
		return nil
	}
	var cfg SMSConfig
	if err := json.Unmarshal([]byte(val.String()), &cfg); err != nil {
		return nil
	}
	return &cfg
}
