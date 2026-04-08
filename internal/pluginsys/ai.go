package pluginsys

import (
	"context"
	"fmt"
)

// ── AI service registration (used by internal/logic/ai) ─────────────────

var _aiServiceFn func(ctx context.Context, action string, params map[string]interface{}) (string, error)

// RegisterAIServiceFn registers the AI call function.
// Called from internal/logic/ai during init().
// Go-native plugins access AI via the PluginContext, not through Goja.
func RegisterAIServiceFn(fn func(ctx context.Context, action string, params map[string]interface{}) (string, error)) {
	_aiServiceFn = fn
}

// CallAIService calls the registered AI service function.
// Exported for use by Go-native plugins via the plugin manager.
func CallAIService(ctx context.Context, action string, params map[string]interface{}) (string, error) {
	if _aiServiceFn == nil {
		return "", fmt.Errorf("AI service not configured")
	}
	return _aiServiceFn(ctx, action, params)
}
