package middleware

import (
	"context"
	"time"

	"github.com/nuxtblog/nuxtblog/internal/langconf"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/i18n/gi18n"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
)

var siteLangCache = gcache.New()

const siteLangCacheTTL = 5 * time.Minute

// InvalidateSiteLangCache clears the cached site default language.
// Call this after updating the site_language option.
func InvalidateSiteLangCache(ctx context.Context) {
	_, _ = siteLangCache.Remove(ctx, "site_language")
}

func getSiteDefaultLang(ctx context.Context) string {
	if v, _ := siteLangCache.Get(ctx, "site_language"); v != nil && v.String() != "" {
		return v.String()
	}
	type OptionRow struct {
		Value string `orm:"value"`
	}
	var row OptionRow
	_ = g.DB().Model("options").Ctx(ctx).Where("key", "site_language").Fields("value").Scan(&row)
	entries := langconf.GetEntries(ctx)
	defaultLocale := "zh-CN"
	if len(entries) > 0 {
		defaultLocale = entries[0].Locale
	}
	lang := defaultLocale
	if row.Value != "" {
		lang = langconf.NormalizeLang(ctx, row.Value)
	}
	_ = siteLangCache.Set(ctx, "site_language", lang, siteLangCacheTTL)
	return lang
}

// LangMiddleware detects the preferred language from the request headers
// and sets it in the context so all g.I18n().T(ctx, key) calls use it.
func LangMiddleware(r *ghttp.Request) {
	lang := detectLang(r)
	r.Request = r.Request.WithContext(gi18n.WithLanguage(r.GetCtx(), lang))
	r.Middleware.Next()
}

func detectLang(r *ghttp.Request) string {
	ctx := r.GetCtx()
	if lang := r.GetHeader("X-Language"); lang != "" {
		return langconf.NormalizeLang(ctx, lang)
	}
	if lang := r.GetHeader("Accept-Language"); lang != "" {
		return langconf.NormalizeLang(ctx, lang)
	}
	return getSiteDefaultLang(ctx)
}
