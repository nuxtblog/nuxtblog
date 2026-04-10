package pluginsys

import (
	"context"
	"fmt"

	sdk "github.com/nuxtblog/nuxtblog/sdk"
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

// ── AI Generate registration (new atomic interface) ─────────────────────

var _aiGenerateFn func(ctx context.Context, req sdk.AIRequest) (*sdk.AIResponse, error)

// RegisterAIGenerateFn registers the Generate implementation.
// Called from internal/logic/ai during init().
func RegisterAIGenerateFn(fn func(ctx context.Context, req sdk.AIRequest) (*sdk.AIResponse, error)) {
	_aiGenerateFn = fn
}

type pluginAIAdapter struct{}

func (pluginAIAdapter) Generate(ctx context.Context, req sdk.AIRequest) (*sdk.AIResponse, error) {
	if _aiGenerateFn == nil {
		return nil, fmt.Errorf("AI service not configured")
	}
	return _aiGenerateFn(ctx, req)
}

// DefaultAI returns the AI adapter that delegates to the registered Generate function.
func DefaultAI() sdk.AI { return pluginAIAdapter{} }
