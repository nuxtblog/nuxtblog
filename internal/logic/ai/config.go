package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/ai/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/google/uuid"
)

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
