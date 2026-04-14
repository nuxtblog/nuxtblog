package ai

import (
	"context"
	"fmt"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	sdk "github.com/nuxtblog/nuxtblog/sdk"
)

type sAI struct{}

func New() service.IAI { return &sAI{} }

func init() {
	service.RegisterAI(New())
	// Legacy AI service function for plugin SDK (nuxtblog.ai) — kept for backward compatibility
	eng.RegisterAIServiceFn(func(ctx context.Context, action string, params map[string]interface{}) (string, error) {
		return "", fmt.Errorf("legacy AI action %q removed; use ctx.AI.Generate() instead", action)
	})

	// Register atomic Generate for plugin SDK (ctx.AI.Generate)
	eng.RegisterAIGenerateFn(func(ctx context.Context, req sdk.AIRequest) (*sdk.AIResponse, error) {
		cfg, err := getActiveConfig(ctx)
		if err != nil {
			return nil, err
		}
		// Convert Messages to system + user strings for current generateText
		var system, user string
		for _, m := range req.Messages {
			switch m.Role {
			case sdk.RoleSystem:
				system = m.Content
			case sdk.RoleUser:
				if user != "" {
					user += "\n\n"
				}
				user += m.Content
			case sdk.RoleAssistant:
				// Best-effort multi-turn with current generateText
				if user != "" {
					user += "\n\nAssistant: " + m.Content + "\n\nUser: "
				}
			}
		}
		text, err := generateText(ctx, *cfg, system, user)
		if err != nil {
			return nil, err
		}
		return &sdk.AIResponse{Text: text}, nil
	})
}

// generateText dispatches to the appropriate API backend based on config format.
func generateText(ctx context.Context, cfg v1.AIConfig, system, user string) (string, error) {
	timeout := time.Duration(cfg.TimeoutMs) * time.Millisecond
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	if cfg.APIFormat == "claude" {
		return claudeGenerate(ctx, cfg, system, user, timeout)
	}
	return openaiGenerate(ctx, cfg, system, user, timeout)
}
