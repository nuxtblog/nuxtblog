package notify

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/nuxtblog/nuxtblog/internal/consts"

	"github.com/gogf/gf/v2/frame/g"
)

func init() {
	Register(&WebhookChannel{})
}

// WebhookConfig is stored in the options table under key "notify_webhook".
type WebhookConfig struct {
	URLs   []string `json:"urls"`   // outbound webhook endpoint URLs
	Secret string   `json:"secret"` // HMAC-SHA256 signing secret; empty = no signature
}

type WebhookChannel struct{}

func (c *WebhookChannel) Name() string { return consts.NotifyChannelWebhook }

func (c *WebhookChannel) Enabled(ctx context.Context) bool {
	cfg := loadWebhookConfig(ctx)
	return cfg != nil && len(cfg.URLs) > 0
}

// Send delivers a user-facing notification event to all configured webhook URLs.
// Satisfies the notify.Channel interface so it participates in notify.Dispatch.
func (c *WebhookChannel) Send(ctx context.Context, msg Message) error {
	cfg := loadWebhookConfig(ctx)
	if cfg == nil || len(cfg.URLs) == 0 {
		return nil
	}
	body, err := json.Marshal(map[string]any{
		"event":        "notification." + msg.Type,
		"type":         msg.Type,
		"sub_type":     msg.SubType,
		"user_id":      msg.UserID,
		"actor_name":   msg.ActorName,
		"object_type":  msg.ObjectType,
		"object_title": msg.ObjectTitle,
		"object_link":  msg.ObjectLink,
		"content":      msg.Content,
		"timestamp":    time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return err
	}
	return dispatchToURLs(ctx, cfg, "notification."+msg.Type, body)
}

// SendWebhookEvent fires a system event to all configured webhook URLs.
// Called directly by listener/webhook.go for non-notification events
// (post.published, comment.created, user.registered, …).
func SendWebhookEvent(ctx context.Context, eventName string, data any) {
	cfg := loadWebhookConfig(ctx)
	if cfg == nil || len(cfg.URLs) == 0 {
		return
	}
	body, err := json.Marshal(map[string]any{
		"event":     eventName,
		"data":      data,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		return
	}
	if err := dispatchToURLs(ctx, cfg, eventName, body); err != nil {
		g.Log().Warningf(ctx, "[webhook] %s error: %v", eventName, err)
	}
}

// dispatchToURLs POSTs body to every URL in cfg. Returns the last error, if any.
func dispatchToURLs(ctx context.Context, cfg *WebhookConfig, eventName string, body []byte) error {
	client := &http.Client{Timeout: 10 * time.Second}
	var lastErr error
	for _, u := range cfg.URLs {
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
		if err != nil {
			lastErr = err
			continue
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Webhook-Event", eventName)
		if cfg.Secret != "" {
			req.Header.Set("X-Webhook-Signature", "sha256="+webhookSignature(body, cfg.Secret))
		}
		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}
		resp.Body.Close()
		if resp.StatusCode >= 400 {
			lastErr = fmt.Errorf("webhook %s returned HTTP %d", u, resp.StatusCode)
		}
	}
	return lastErr
}

func webhookSignature(body []byte, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}

func loadWebhookConfig(ctx context.Context) *WebhookConfig {
	val, err := g.DB().Ctx(ctx).Model("options").Where("key", "notify_webhook").Value("value")
	if err != nil || val.IsEmpty() {
		return nil
	}
	var cfg WebhookConfig
	if err := json.Unmarshal([]byte(val.String()), &cfg); err != nil {
		return nil
	}
	return &cfg
}
