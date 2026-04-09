package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v3"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// ── Marketplace ───────────────────────────────────────────────────────────

const (
	defaultRegistryURL  = "https://raw.githubusercontent.com/nuxtblog/registry/main/registry.json"
	registryCacheKey    = "plugin_registry"
	registrySyncedAtKey = "plugin_registry_synced_at"
	registryTTL         = time.Hour
)

// GetMarketplace returns the marketplace plugin list, auto-syncing if the cache is stale.
func (s *sPlugin) GetMarketplace(ctx context.Context, keyword, pluginType string) (*v1.MarketplaceListRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	syncedAt := getOption(ctx, registrySyncedAtKey)
	needsSync := syncedAt == ""
	if !needsSync {
		t, err := time.Parse(time.RFC3339, syncedAt)
		needsSync = err != nil || time.Since(t) > registryTTL
	}
	if needsSync {
		if res, err := s.SyncMarketplace(ctx); err == nil {
			syncedAt = res.SyncedAt
		}
	}

	var items []v1.MarketplaceItem
	if cache := getOption(ctx, registryCacheKey); cache != "" {
		_ = json.Unmarshal([]byte(cache), &items)
	}

	if keyword != "" || pluginType != "" {
		q := strings.ToLower(keyword)
		out := items[:0]
		for _, item := range items {
			if pluginType != "" && item.Type != pluginType {
				continue
			}
			if q != "" &&
				!strings.Contains(strings.ToLower(item.Title), q) &&
				!strings.Contains(strings.ToLower(item.Description), q) &&
				!strings.Contains(strings.ToLower(item.Author), q) {
				continue
			}
			out = append(out, item)
		}
		items = out
	}

	if items == nil {
		items = []v1.MarketplaceItem{}
	}
	return &v1.MarketplaceListRes{Items: items, SyncedAt: syncedAt}, nil
}

// SyncMarketplace fetches the latest registry.json and stores it in the options table.
func (s *sPlugin) SyncMarketplace(ctx context.Context) (*v1.MarketplaceSyncRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	regURL := g.Cfg().MustGet(ctx, "plugin.registry_url", defaultRegistryURL).String()
	body, err := pluginHTTPGet(ctx, regURL, true)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "fetch registry failed: "+err.Error())
	}

	var items []v1.MarketplaceItem
	if err := json.Unmarshal(body, &items); err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError, "invalid registry format: "+err.Error())
	}

	now := time.Now().UTC().Format(time.RFC3339)
	setOption(ctx, registryCacheKey, string(body))
	setOption(ctx, registrySyncedAtKey, now)

	return &v1.MarketplaceSyncRes{Count: len(items), SyncedAt: now}, nil
}

func getOption(ctx context.Context, key string) string {
	type row struct {
		Value string `orm:"value"`
	}
	var r row
	_ = g.DB().Ctx(ctx).Model("options").Where("key", key).Fields("value").Scan(&r)
	return r.Value
}

func setOption(ctx context.Context, key, value string) {
	cnt, _ := g.DB().Ctx(ctx).Model("options").Where("key", key).Count()
	if cnt > 0 {
		_, _ = g.DB().Ctx(ctx).Model("options").Where("key", key).
			Data(g.Map{"value": value}).Update()
	} else {
		_, _ = g.DB().Ctx(ctx).Model("options").
			Data(g.Map{"key": key, "value": value}).Insert()
	}
}

// ── GitHub proxy helpers ──────────────────────────────────────────────────

// applyGitHubProxy rewrites a GitHub URL using the configured mirror prefix.
// e.g. "https://api.github.com/..." → "https://ghproxy.net/https://api.github.com/..."
func applyGitHubProxy(ctx context.Context, rawURL string) string {
	proxy := strings.TrimRight(getOption(ctx, "plugin_github_proxy"), "/")
	if proxy == "" {
		return rawURL
	}
	for _, prefix := range []string{
		"https://api.github.com",
		"https://github.com",
		"https://raw.githubusercontent.com",
		"https://objects.githubusercontent.com",
	} {
		if strings.HasPrefix(rawURL, prefix) {
			return proxy + "/" + rawURL
		}
	}
	return rawURL
}

// pluginHTTPGet performs an HTTP GET like githubGet, but honours the
// plugin_http_proxy and plugin_github_proxy options.
func pluginHTTPGet(ctx context.Context, rawURL string, noCache ...bool) ([]byte, error) {
	targetURL := applyGitHubProxy(ctx, rawURL)

	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxyAddr := getOption(ctx, "plugin_http_proxy"); proxyAddr != "" {
		proxyURL, err := url.Parse(proxyAddr)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}
	client := &http.Client{Timeout: 60 * time.Second, Transport: transport}

	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "nuxtblog-plugin-installer/1.0")
	req.Header.Set("Accept", "application/vnd.github+json")
	if len(noCache) > 0 && noCache[0] {
		req.Header.Set("Cache-Control", "no-cache")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("not found (404)")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// ── Preview ───────────────────────────────────────────────────────────────

// previewCache is a process-local TTL cache for plugin preview responses,
// keyed by "owner/repo". Preview data is ephemeral; no need to persist across
// restarts, so an in-memory sync.Map is faster and simpler than the options
// table used by the marketplace cache.
type previewCacheEntry struct {
	res *v1.PluginPreviewRes
	at  time.Time
}

var (
	previewCache    sync.Map
	previewCacheTTL = 10 * time.Minute
)

// Preview fetches plugin.yaml (or package.json) from the GitHub repo's default branch
// and returns parsed manifest metadata without downloading or installing the plugin.
func (s *sPlugin) Preview(ctx context.Context, repo string) (*v1.PluginPreviewRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	owner, repoName, err := parseRepo(repo)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, err.Error())
	}
	cacheKey := owner + "/" + repoName

	// Serve from cache if fresh
	if v, ok := previewCache.Load(cacheKey); ok {
		e := v.(*previewCacheEntry)
		if time.Since(e.at) < previewCacheTTL {
			return e.res, nil
		}
		previewCache.Delete(cacheKey)
	}

	// Try plugin.yaml first
	yamlURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/HEAD/plugin.yaml", owner, repoName)
	if yamlData, yamlErr := pluginHTTPGet(ctx, yamlURL); yamlErr == nil {
		mf, _, parseErr := parsePluginYAML(yamlData)
		if parseErr == nil {
			res := buildPreviewFromManifest(mf)
			// Parse depends from raw YAML (eng.Manifest doesn't carry depends)
			var raw pluginYAMLManifest
			if yaml.Unmarshal(yamlData, &raw) == nil && len(raw.Depends) > 0 {
				for _, d := range raw.Depends {
					res.Depends = append(res.Depends, v1.DependencyPreview{
						ID: d.ID, Version: d.Version, Optional: d.Optional,
					})
				}
			}
			previewCache.Store(cacheKey, &previewCacheEntry{res: res, at: time.Now()})
			return res, nil
		}
	}

	// Fallback: package.json
	rawURL := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/HEAD/package.json", owner, repoName)
	data, err := pluginHTTPGet(ctx, rawURL)
	if err != nil {
		return nil, gerror.NewCode(gcode.CodeInternalError,
			fmt.Sprintf("cannot fetch plugin.yaml or package.json from %s/%s: %v", owner, repoName, err))
	}

	var pkg struct {
		Name        string `json:"name"`
		Version     string `json:"version"`
		Description string `json:"description"`
		Author      any    `json:"author"`
		Plugin      *struct {
			Title        string                 `json:"title"`
			Icon         string                 `json:"icon"`
			CSS          string                 `json:"css"`
			Priority     int                    `json:"priority"`
			Settings     []eng.SettingField     `json:"settings"`
			Capabilities eng.PluginCapabilities `json:"capabilities"`
			Webhooks     []eng.WebhookDef       `json:"webhooks"`
			Pipelines    []eng.PipelineDef      `json:"pipelines"`
		} `json:"plugin"`
	}
	if err = json.Unmarshal(data, &pkg); err != nil || pkg.Plugin == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter,
			"package.json must contain a \"plugin\" field")
	}

	author := ""
	switch v := pkg.Author.(type) {
	case string:
		author = v
	case map[string]any:
		if s, ok2 := v["name"].(string); ok2 {
			author = s
		}
	}

	settings := make([]v1.SettingField, 0, len(pkg.Plugin.Settings))
	for _, sf := range pkg.Plugin.Settings {
		settings = append(settings, v1.SettingField{
			Key:         sf.Key,
			Label:       sf.Label,
			Type:        string(sf.Type),
			Required:    sf.Required,
			Default:     sf.Default,
			Placeholder: sf.Placeholder,
			Description: sf.Description,
			Options:     sf.Options,
			Shared:      sf.Shared,
		})
	}

	webhooks := make([]v1.WebhookPreview, 0, len(pkg.Plugin.Webhooks))
	for _, wh := range pkg.Plugin.Webhooks {
		webhooks = append(webhooks, v1.WebhookPreview{URL: wh.URL, Events: wh.Events})
	}

	pipelines := make([]v1.PipelinePreview, 0, len(pkg.Plugin.Pipelines))
	for _, p := range pkg.Plugin.Pipelines {
		pipelines = append(pipelines, v1.PipelinePreview{
			Name:      p.Name,
			Trigger:   p.Trigger,
			StepCount: len(p.Steps),
		})
	}

	caps := v1.CapabilitiesPreview{}
	if pkg.Plugin.Capabilities.HTTP != nil {
		caps.HTTP = &v1.HTTPCapPreview{
			Allow:     pkg.Plugin.Capabilities.HTTP.Allow,
			TimeoutMs: pkg.Plugin.Capabilities.HTTP.TimeoutMs,
		}
	}
	if pkg.Plugin.Capabilities.Store != nil {
		caps.Store = &v1.StoreCapPreview{
			Read:  pkg.Plugin.Capabilities.Store.Read,
			Write: pkg.Plugin.Capabilities.Store.Write,
		}
	}
	if pkg.Plugin.Capabilities.Events != nil {
		caps.Events = &v1.EventsCapPreview{
			Subscribe: pkg.Plugin.Capabilities.Events.Subscribe,
		}
	}
	if pkg.Plugin.Capabilities.DB != nil {
		caps.DB = convertDBCapPreview(pkg.Plugin.Capabilities.DB)
	}

	res := &v1.PluginPreviewRes{
		Name:         pkg.Name,
		Title:        pkg.Plugin.Title,
		Description:  pkg.Description,
		Version:      pkg.Version,
		Author:       author,
		Icon:         orDefault(pkg.Plugin.Icon, "i-tabler-plug"),
		Priority:     pkg.Plugin.Priority,
		HasCSS:       pkg.Plugin.CSS != "",
		Capabilities: caps,
		Settings:     settings,
		Webhooks:     webhooks,
		Pipelines:    pipelines,
	}
	previewCache.Store(cacheKey, &previewCacheEntry{res: res, at: time.Now()})
	return res, nil
}

// buildPreviewFromManifest converts an engine Manifest into a PluginPreviewRes.
func buildPreviewFromManifest(mf *eng.Manifest) *v1.PluginPreviewRes {
	settings := make([]v1.SettingField, 0, len(mf.Settings))
	for _, sf := range mf.Settings {
		settings = append(settings, v1.SettingField{
			Key:         sf.Key,
			Label:       sf.Label,
			Type:        string(sf.Type),
			Required:    sf.Required,
			Default:     sf.Default,
			Placeholder: sf.Placeholder,
			Description: sf.Description,
			Options:     sf.Options,
			Shared:      sf.Shared,
		})
	}

	caps := v1.CapabilitiesPreview{}
	if mf.Capabilities.DB != nil {
		caps.DB = convertDBCapPreview(mf.Capabilities.DB)
	}

	return &v1.PluginPreviewRes{
		Name:         mf.Name,
		Title:        mf.Title,
		Description:  mf.Description,
		Version:      mf.Version,
		Author:       mf.Author,
		Icon:         orDefault(mf.Icon, "i-tabler-plug"),
		HasCSS:       mf.CSS != "",
		Capabilities: caps,
		Settings:     settings,
	}
}

// convertDBCapPreview converts engine DBCap to API DBCapPreview.
func convertDBCapPreview(d *eng.DBCap) *v1.DBCapPreview {
	if d == nil {
		return nil
	}
	r := &v1.DBCapPreview{
		Own: d.Own,
		Raw: d.Raw,
	}
	for _, t := range d.Tables {
		r.Tables = append(r.Tables, v1.DBTablePreview{
			Table: t.Table,
			Ops:   t.Ops,
		})
	}
	return r
}
