package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	sdk "github.com/nuxtblog/nuxtblog/sdk"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
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

// ── Options keys ──────────────────────────────────────────────────────────────

const (
	optKeyConfigs  = "ai_configs"
	optKeyActiveID = "ai_active_id"
)

// ── Config storage helpers ────────────────────────────────────────────────────

func loadConfigs(ctx context.Context) ([]v1.AIConfig, error) {
	type row struct{ Value string `orm:"value"` }
	var r row
	_ = g.DB().Ctx(ctx).Model("options").Where("key", optKeyConfigs).Scan(&r)
	if r.Value == "" {
		return nil, nil
	}
	var configs []v1.AIConfig
	if err := json.Unmarshal([]byte(r.Value), &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func saveConfigs(ctx context.Context, configs []v1.AIConfig) error {
	b, err := json.Marshal(configs)
	if err != nil {
		return err
	}
	val := string(b)
	cnt, _ := g.DB().Ctx(ctx).Model("options").Where("key", optKeyConfigs).Count()
	if cnt > 0 {
		_, err = g.DB().Ctx(ctx).Model("options").Where("key", optKeyConfigs).
			Data(g.Map{"value": val}).Update()
	} else {
		_, err = g.DB().Ctx(ctx).Model("options").
			Data(g.Map{"key": optKeyConfigs, "value": val, "autoload": 0}).Insert()
	}
	return err
}

func loadActiveID(ctx context.Context) string {
	type row struct{ Value string `orm:"value"` }
	var r row
	_ = g.DB().Ctx(ctx).Model("options").Where("key", optKeyActiveID).Scan(&r)
	return r.Value
}

func saveActiveID(ctx context.Context, id string) error {
	cnt, _ := g.DB().Ctx(ctx).Model("options").Where("key", optKeyActiveID).Count()
	var err error
	if cnt > 0 {
		_, err = g.DB().Ctx(ctx).Model("options").Where("key", optKeyActiveID).
			Data(g.Map{"value": id}).Update()
	} else {
		_, err = g.DB().Ctx(ctx).Model("options").
			Data(g.Map{"key": optKeyActiveID, "value": id, "autoload": 0}).Insert()
	}
	return err
}

// getActiveConfig returns the active config with the real (unmasked) API key.
func getActiveConfig(ctx context.Context) (*v1.AIConfig, error) {
	configs, err := loadConfigs(ctx)
	if err != nil {
		return nil, err
	}
	activeID := loadActiveID(ctx)
	if activeID == "" && len(configs) > 0 {
		activeID = configs[0].ID
	}
	for _, c := range configs {
		if c.ID == activeID {
			return &c, nil
		}
	}
	return nil, gerror.NewCode(gcode.CodeInvalidOperation,
		"no active AI config — please configure one in Admin → AI")
}

func getConfigByID(ctx context.Context, id string) (*v1.AIConfig, error) {
	configs, err := loadConfigs(ctx)
	if err != nil {
		return nil, err
	}
	for _, c := range configs {
		if c.ID == id {
			return &c, nil
		}
	}
	return nil, gerror.NewCode(gcode.CodeNotFound, "config not found")
}

func maskAPIKey(key string) string {
	if len(key) > 8 {
		return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
	}
	if key != "" {
		return "****"
	}
	return ""
}

func requireAdmin(ctx context.Context) error {
	role := middleware.GetCurrentUserRole(ctx)
	if role < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized, "admin required")
	}
	return nil
}

// ── Config CRUD ───────────────────────────────────────────────────────────────

func (s *sAI) ListConfigs(ctx context.Context) (*v1.AIListConfigsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	configs, err := loadConfigs(ctx)
	if err != nil {
		return nil, err
	}
	if configs == nil {
		configs = []v1.AIConfig{}
	}
	activeID := loadActiveID(ctx)
	for i := range configs {
		configs[i].IsActive = configs[i].ID == activeID
		configs[i].APIKey = maskAPIKey(configs[i].APIKey)
	}
	return &v1.AIListConfigsRes{Items: configs, ActiveID: activeID}, nil
}

func (s *sAI) CreateConfig(ctx context.Context, req *v1.AICreateConfigReq) (*v1.AICreateConfigRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	configs, err := loadConfigs(ctx)
	if err != nil {
		return nil, err
	}
	item := v1.AIConfig{
		ID:        uuid.New().String()[:8],
		Name:      req.Name,
		APIFormat: req.APIFormat,
		Label:     req.Label,
		APIKey:    req.APIKey,
		Model:     req.Model,
		BaseURL:   req.BaseURL,
		TimeoutMs: req.TimeoutMs,
	}
	if item.TimeoutMs <= 0 {
		item.TimeoutMs = 30000
	}
	configs = append(configs, item)
	if err := saveConfigs(ctx, configs); err != nil {
		return nil, err
	}
	// Auto-activate first config
	if len(configs) == 1 {
		_ = saveActiveID(ctx, item.ID)
		item.IsActive = true
	}
	ret := item
	ret.APIKey = maskAPIKey(ret.APIKey)
	return &v1.AICreateConfigRes{Item: ret}, nil
}

func (s *sAI) UpdateConfig(ctx context.Context, req *v1.AIUpdateConfigReq) (*v1.AIUpdateConfigRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	configs, err := loadConfigs(ctx)
	if err != nil {
		return nil, err
	}
	found := false
	var updated v1.AIConfig
	for i := range configs {
		if configs[i].ID != req.ID {
			continue
		}
		if req.Name != "" {
			configs[i].Name = req.Name
		}
		if req.APIFormat != "" {
			configs[i].APIFormat = req.APIFormat
		}
		configs[i].Label = req.Label
		// Only update API key if it's a real new value (not a masked placeholder)
		if req.APIKey != "" && !strings.Contains(req.APIKey, "****") {
			configs[i].APIKey = req.APIKey
		}
		if req.Model != "" {
			configs[i].Model = req.Model
		}
		configs[i].BaseURL = req.BaseURL
		if req.TimeoutMs > 0 {
			configs[i].TimeoutMs = req.TimeoutMs
		}
		updated = configs[i]
		found = true
		break
	}
	if !found {
		return nil, gerror.NewCode(gcode.CodeNotFound, "config not found")
	}
	if err := saveConfigs(ctx, configs); err != nil {
		return nil, err
	}
	updated.APIKey = maskAPIKey(updated.APIKey)
	return &v1.AIUpdateConfigRes{Item: updated}, nil
}

func (s *sAI) DeleteConfig(ctx context.Context, id string) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	configs, err := loadConfigs(ctx)
	if err != nil {
		return err
	}
	filtered := make([]v1.AIConfig, 0, len(configs))
	for _, c := range configs {
		if c.ID != id {
			filtered = append(filtered, c)
		}
	}
	if err := saveConfigs(ctx, filtered); err != nil {
		return err
	}
	if loadActiveID(ctx) == id {
		newActive := ""
		if len(filtered) > 0 {
			newActive = filtered[0].ID
		}
		_ = saveActiveID(ctx, newActive)
	}
	return nil
}

func (s *sAI) ActivateConfig(ctx context.Context, id string) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	configs, err := loadConfigs(ctx)
	if err != nil {
		return err
	}
	for _, c := range configs {
		if c.ID == id {
			return saveActiveID(ctx, id)
		}
	}
	return gerror.NewCode(gcode.CodeNotFound, "config not found")
}

func (s *sAI) TestConfig(ctx context.Context, id string) (*v1.AITestConfigRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	cfg, err := getConfigByID(ctx, id)
	if err != nil {
		return nil, err
	}
	result, testErr := generateText(ctx, *cfg, "", "Reply with exactly one word: OK")
	if testErr != nil {
		return &v1.AITestConfigRes{OK: false, Message: testErr.Error()}, nil
	}
	msg := "Connection successful"
	if result != "" {
		msg = fmt.Sprintf("Connection successful — model replied: %s", strings.TrimSpace(result))
	}
	return &v1.AITestConfigRes{OK: true, Message: msg}, nil
}

func (s *sAI) FromURL(ctx context.Context, url, style string) (*v1.AIFromURLRes, error) {
	cfg, err := getActiveConfig(ctx)
	if err != nil {
		return nil, err
	}
	// Fetch the page content
	pageText, fetchErr := fetchPageText(ctx, url)
	if fetchErr != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidOperation,
			fmt.Sprintf("fetch URL failed: %v", fetchErr))
	}
	// Truncate to 8000 chars to fit in context
	if len(pageText) > 8000 {
		pageText = pageText[:8000] + "..."
	}
	stylePart := ""
	switch style {
	case "formal":
		stylePart = "语气正式专业。"
	case "casual":
		stylePart = "语气轻松口语。"
	case "concise":
		stylePart = "简洁为主。"
	}
	system := fmt.Sprintf(`你是一位博客写手。请根据以下网页内容改写成一篇原创博客文章。%s
请以如下 JSON 格式输出（不要添加 markdown 代码块）：
{"title":"...","content":"...","excerpt":"..."}
其中 content 为 HTML 格式正文。`, stylePart)

	raw, err := generateText(ctx, *cfg, system, pageText)
	if err != nil {
		return nil, err
	}
	// Parse JSON response
	raw = strings.TrimSpace(raw)
	// Strip markdown code fences if model added them
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	var res v1.AIFromURLRes
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		// Fallback: return raw as content
		res.Title = "从 URL 生成的文章"
		res.Content = raw
	}
	return &res, nil
}

// ── HTTP helper: fetch page text ──────────────────────────────────────────────

func fetchPageText(ctx context.Context, url string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NuxtBlog/2.0)")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 512*1024)) // max 512 KB
	if err != nil {
		return "", err
	}
	// Very basic HTML → text: strip tags
	text := string(body)
	// Remove script/style blocks
	for _, tag := range []string{"script", "style", "head"} {
		for {
			start := strings.Index(strings.ToLower(text), "<"+tag)
			end := strings.Index(strings.ToLower(text), "</"+tag+">")
			if start < 0 || end < 0 || end < start {
				break
			}
			text = text[:start] + text[end+len("</"+tag+">"):]
		}
	}
	// Strip remaining HTML tags
	var b strings.Builder
	inTag := false
	for _, c := range text {
		if c == '<' {
			inTag = true
		} else if c == '>' {
			inTag = false
			b.WriteRune(' ')
		} else if !inTag {
			b.WriteRune(c)
		}
	}
	// Collapse whitespace
	result := b.String()
	lines := strings.Split(result, "\n")
	var clean []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			clean = append(clean, line)
		}
	}
	return strings.Join(clean, "\n"), nil
}

// ── AI call: unified interface ────────────────────────────────────────────────

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

// ── OpenAI-compatible call ────────────────────────────────────────────────────

type oaiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type oaiRequest struct {
	Model       string       `json:"model"`
	Messages    []oaiMessage `json:"messages"`
	MaxTokens   int          `json:"max_tokens,omitempty"`
	Temperature float64      `json:"temperature"`
}

type oaiResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func openaiGenerate(ctx context.Context, cfg v1.AIConfig, system, user string, timeout time.Duration) (string, error) {
	base := cfg.BaseURL
	if base == "" {
		return "", gerror.New("base_url is required — please set it in Admin → AI config")
	}
	base = strings.TrimRight(base, "/")
	endpoint := base + "/v1/chat/completions"

	msgs := []oaiMessage{}
	if system != "" {
		msgs = append(msgs, oaiMessage{Role: "system", Content: system})
	}
	msgs = append(msgs, oaiMessage{Role: "user", Content: user})

	body, _ := json.Marshal(oaiRequest{
		Model:       cfg.Model,
		Messages:    msgs,
		Temperature: 0.7,
	})

	httpCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var oaiRes oaiResponse
	if err := json.Unmarshal(data, &oaiRes); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if oaiRes.Error != nil {
		return "", gerror.New(oaiRes.Error.Message)
	}
	if len(oaiRes.Choices) == 0 {
		return "", gerror.New("AI returned no content")
	}
	return oaiRes.Choices[0].Message.Content, nil
}

// ── Claude (Anthropic) call ───────────────────────────────────────────────────

type claudeRequest struct {
	Model     string       `json:"model"`
	MaxTokens int          `json:"max_tokens"`
	System    string       `json:"system,omitempty"`
	Messages  []oaiMessage `json:"messages"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func claudeGenerate(ctx context.Context, cfg v1.AIConfig, system, user string, timeout time.Duration) (string, error) {
	base := cfg.BaseURL
	if base == "" {
		base = "https://api.anthropic.com"
	}
	base = strings.TrimRight(base, "/")
	endpoint := base + "/v1/messages"

	payload := claudeRequest{
		Model:     cfg.Model,
		MaxTokens: 4096,
		System:    system,
		Messages:  []oaiMessage{{Role: "user", Content: user}},
	}
	body, _ := json.Marshal(payload)

	httpCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", cfg.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var claudeRes claudeResponse
	if err := json.Unmarshal(data, &claudeRes); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	if claudeRes.Error != nil {
		return "", gerror.New(claudeRes.Error.Message)
	}
	if len(claudeRes.Content) == 0 {
		return "", gerror.New("AI returned no content")
	}
	return claudeRes.Content[0].Text, nil
}
