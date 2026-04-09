package plugin

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type sPlugin struct{}

func New() service.IPlugin { return &sPlugin{} }
func init()                { service.RegisterPlugin(New()) }

// ── List ──────────────────────────────────────────────────────────────────

func (s *sPlugin) List(ctx context.Context) (*v1.PluginListRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type row struct {
		Id           string `orm:"id"`
		Title        string `orm:"title"`
		Description  string `orm:"description"`
		Version      string `orm:"version"`
		Author       string `orm:"author"`
		Icon         string `orm:"icon"`
		RepoUrl      string `orm:"repo_url"`
		Enabled      int    `orm:"enabled"`
		InstalledAt  string `orm:"installed_at"`
		Capabilities string `orm:"capabilities"`
		Source       string `orm:"source"`
		Manifest     string `orm:"manifest"`
	}
	var rows []row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Fields("id,title,description,version,author,icon,repo_url,enabled,installed_at,COALESCE(capabilities,'{}') as capabilities,COALESCE(source,'external') as source,COALESCE(manifest,'{}') as manifest").
		OrderDesc("installed_at").
		Scan(&rows); err != nil {
		return nil, err
	}

	items := make([]v1.PluginItem, 0, len(rows))
	for _, r := range rows {
		source := r.Source
		if source == "" {
			source = "external"
		}
		// Extract type from stored manifest JSON
		pluginType := ""
		if r.Manifest != "" {
			var mf struct {
				Type string `json:"type"`
			}
			if json.Unmarshal([]byte(r.Manifest), &mf) == nil {
				pluginType = mf.Type
			}
		}
		// Derive type from source if not in manifest
		if pluginType == "" {
			if source == "builtin" {
				pluginType = "builtin"
			} else {
				pluginType = "js" // default for external
			}
		}
		items = append(items, v1.PluginItem{
			Id:           r.Id,
			Title:        r.Title,
			Description:  r.Description,
			Version:      r.Version,
			Author:       r.Author,
			Icon:         r.Icon,
			RepoUrl:      r.RepoUrl,
			Enabled:      r.Enabled == 1,
			InstalledAt:  r.InstalledAt,
			Capabilities: r.Capabilities,
			Source:       source,
			Type:         pluginType,
		})
	}

	// Merge Layer 0 (YAML declarative plugins) — always enabled, no DB record
	for _, yp := range eng.GetAllYAMLPlugins() {
		items = append(items, v1.PluginItem{
			Id:          yp.ID,
			Title:       yp.Title,
			Description: yp.Description,
			Version:     yp.Version,
			Author:      yp.Author,
			Icon:        yp.Icon,
			Enabled:     true,
			Type:        "yaml",
		})
	}

	return &v1.PluginListRes{Items: items}, nil
}

// ── helpers ───────────────────────────────────────────────────────────────

func requireAdmin(ctx context.Context) error {
	if middleware.GetCurrentUserRole(ctx) < middleware.RoleAdmin {
		return gerror.NewCode(gcode.CodeNotAuthorized,
			g.I18n().T(ctx, "error.forbidden"))
	}
	return nil
}

// parseRepo normalises various GitHub URL formats to (owner, repo).
func parseRepo(raw string) (string, string, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")
	raw = strings.TrimPrefix(raw, "github.com/")
	parts := strings.SplitN(raw, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo format, expected owner/repo or github.com/owner/repo")
	}
	// Strip .git suffix and any trailing path
	repo := strings.TrimSuffix(strings.Split(parts[1], "/")[0], ".git")
	return parts[0], repo, nil
}

var (
	httpClient   = &http.Client{Timeout: 20 * time.Second}
	githubClient = &http.Client{Timeout: 60 * time.Second}
)

// httpGet is a plain HTTP GET used for non-GitHub URLs (e.g. registry sync).
func httpGet(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}
	return io.ReadAll(resp.Body)
}

// githubGet makes a GitHub API / CDN request with the required User-Agent header.
// It follows redirects (needed for zipball URLs) and uses a longer timeout.
func githubGet(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "nuxtblog-plugin-installer/1.0")
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := githubClient.Do(req)
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

func orDefault(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

// manifestSettingsJSON serialises the plugin's settings schema to a JSON string.
func manifestSettingsJSON(m *eng.Manifest) string {
	if len(m.Settings) == 0 {
		return "[]"
	}
	b, _ := json.Marshal(m.Settings)
	return string(b)
}

// manifestDefaultSettings builds a JSON object pre-filled with each field's
// default value so plugins can call blog.settings.get(key) immediately.
func manifestDefaultSettings(m *eng.Manifest) string {
	defaults := make(map[string]interface{}, len(m.Settings))
	for _, f := range m.Settings {
		if f.Default != nil {
			defaults[f.Key] = f.Default
		}
	}
	b, _ := json.Marshal(defaults)
	return string(b)
}

// manifestCapabilitiesJSON serialises the plugin's capability declarations.
func manifestCapabilitiesJSON(m *eng.Manifest) string {
	b, _ := json.Marshal(m.Capabilities)
	if b == nil {
		return "{}"
	}
	return string(b)
}

// manifestFullJSON serialises the complete plugin manifest to a JSON string.
func manifestFullJSON(m *eng.Manifest) string {
	b, err := json.Marshal(m)
	if err != nil {
		return "{}"
	}
	return string(b)
}
