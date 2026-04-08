package listener

import (
	"context"

	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/notify"
)

// registerWebhookListeners subscribes to key system events and forwards them
// to all configured webhook URLs (see notify_webhook option in admin settings).
//
// Payload is serialised as-is, so receivers (n8n, Zapier, etc.) get a typed
// JSON body: { "event": "post.published", "data": { ... }, "timestamp": "..." }
func registerWebhookListeners() {
	webhookEvents := []string{
		event.PostPublished,
		event.PostCreated,
		event.PostUpdated,
		event.PostDeleted,
		event.CommentCreated,
		event.UserRegistered,
		event.MediaUploaded,
	}
	for _, name := range webhookEvents {
		n := name // capture loop variable
		event.OnAsync(n, func(ctx context.Context, e event.Event) error {
			notify.SendWebhookEvent(ctx, e.Name, e.Payload)
			return nil
		})
	}
}
