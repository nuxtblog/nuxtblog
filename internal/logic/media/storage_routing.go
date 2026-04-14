package media

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/storage"
)

// ── Upload path ───────────────────────────────────────────────────────────────

func getUploadPathTemplate(ctx context.Context) string {
	val, err := dao.Options.Ctx(ctx).Where("key", "media_upload_path").Value("value")
	if err != nil || val.IsNil() {
		return consts.StoragePathTemplateYearMonth
	}
	raw := val.String()
	var s string
	if json.Unmarshal([]byte(raw), &s) == nil && s != "" {
		return s
	}
	if raw != "" {
		return raw
	}
	return consts.StoragePathTemplateYearMonth
}

// getCategoryPathTemplate returns the path template for a category, falling back to global.
func getCategoryPathTemplate(ctx context.Context, category string) string {
	cat := GetCategoryDef(ctx, category)
	if cat != nil && cat.PathTemplate != "" {
		return cat.PathTemplate
	}
	return getUploadPathTemplate(ctx)
}

// ── Storage routing ───────────────────────────────────────────────────────────

type storageRule struct {
	MimePrefix string `json:"mimePrefix"`
	Backend    string `json:"backend"`
}

type storageRouting struct {
	Default string        `json:"default"`
	Rules   []storageRule `json:"rules"`
}

// resolveBackend picks the right storage backend.
// Priority: category storage key (from media_categories option) > MIME prefix rule > default.
func resolveBackend(ctx context.Context, mimeType, category string) (backendName string, up storage.Uploader) {
	// 1. Category storage key from media_categories option (takes highest priority)
	if category != "" {
		if key := GetCategoryStorageKey(ctx, category); key != "" {
			return key, storage.Named(ctx, key)
		}
	}

	// 2. MIME prefix rules from storage_routing option
	var routing storageRouting
	_ = getOptionJSON(ctx, "storage_routing", &routing)
	for _, rule := range routing.Rules {
		if rule.MimePrefix != "" && strings.HasPrefix(mimeType, rule.MimePrefix) {
			return rule.Backend, storage.Named(ctx, rule.Backend)
		}
	}

	// 3. Routing default, then config default
	if routing.Default != "" {
		return routing.Default, storage.Named(ctx, routing.Default)
	}
	name := storage.DefaultName(ctx)
	return name, storage.Named(ctx, name)
}
