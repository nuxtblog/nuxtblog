// Package langconf provides helpers for reading the supported-language list
// from config (i18n.languages). Both the lang middleware and the Languages API
// use this so that adding a new language only requires a config change.
package langconf

import (
	"context"
	"strings"

	v1 "github.com/nuxtblog/nuxtblog/api/site/v1"

	"github.com/gogf/gf/v2/frame/g"
)

// Entry mirrors one item in the i18n.languages config array.
type Entry struct {
	Code   string `json:"code"`   // API code, e.g. "zh"
	Name   string `json:"name"`   // native name, e.g. "中文"
	Label  string `json:"label"`  // English label, e.g. "Chinese"
	Prefix string `json:"prefix"` // Accept-Language prefix to match, e.g. "zh"
	Locale string `json:"locale"` // gi18n locale string, e.g. "zh-CN"
}

// fallback is used when the config key is absent or empty.
var fallback = []Entry{
	{Code: "zh", Name: "中文", Label: "Chinese", Prefix: "zh", Locale: "zh-CN"},
	{Code: "en", Name: "English", Label: "English", Prefix: "en", Locale: "en"},
}

// GetEntries returns the configured language list, falling back to the built-in
// defaults when the config key is missing.
func GetEntries(ctx context.Context) []Entry {
	val, err := g.Cfg().Get(ctx, "i18n.languages")
	if err != nil || val.IsNil() {
		return fallback
	}
	var entries []Entry
	if err := val.Scan(&entries); err != nil || len(entries) == 0 {
		return fallback
	}
	return entries
}

// GetLanguages returns the language list in the API response shape.
func GetLanguages(ctx context.Context) []v1.Language {
	entries := GetEntries(ctx)
	langs := make([]v1.Language, len(entries))
	for i, e := range entries {
		langs[i] = v1.Language{Code: e.Code, Name: e.Name, Label: e.Label}
	}
	return langs
}

// NormalizeLang maps a raw Accept-Language / X-Language header value to the
// canonical locale string used by gi18n (e.g. "zh-CN", "en").
// Falls back to the first configured language's locale.
func NormalizeLang(ctx context.Context, lang string) string {
	entries := GetEntries(ctx)
	parts := strings.Split(lang, ",")
	first := strings.ToLower(strings.TrimSpace(strings.Split(parts[0], ";")[0]))
	for _, e := range entries {
		if strings.HasPrefix(first, e.Prefix) {
			return e.Locale
		}
	}
	if len(entries) > 0 {
		return entries[0].Locale
	}
	return "zh-CN"
}
