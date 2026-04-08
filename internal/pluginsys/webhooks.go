package pluginsys

// webhooks.go — standalone webhook utilities for YAML and Go-native plugins.
//
// The former Goja runtime-bound webhook methods have been removed.

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"time"

	"github.com/gogf/gf/v2/frame/g"
)

const webhookHTTPTimeout = 10 * time.Second

var webhookHTTPClient = &http.Client{Timeout: webhookHTTPTimeout}

// settingsRef matches {{settings.key}} placeholders in webhook URL / header values.
var settingsRef = regexp.MustCompile(`\{\{settings\.([^}]+)\}\}`)

// fireWebhookHTTP is a standalone webhook sender for YAML plugins (no runtime needed).
func fireWebhookHTTP(url string, headers map[string]string, eventName string, payload map[string]any) {
	body, err := json.Marshal(payload)
	if err != nil {
		g.Log().Warningf(context.Background(), "[yaml-plugin] webhook marshal: %v", err)
		return
	}
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		g.Log().Warningf(context.Background(), "[yaml-plugin] webhook request %s: %v", url, err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	resp, err := webhookHTTPClient.Do(req)
	if err != nil {
		g.Log().Warningf(context.Background(), "[yaml-plugin] webhook POST %s: %v", url, err)
		return
	}
	defer func() { _, _ = io.Copy(io.Discard, resp.Body); resp.Body.Close() }()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		g.Log().Warningf(context.Background(), "[yaml-plugin] webhook POST %s → HTTP %d", url, resp.StatusCode)
	}
}
