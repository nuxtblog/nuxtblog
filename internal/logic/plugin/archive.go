package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
	"gopkg.in/yaml.v3"

	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
)

// isSkippedArchivePath returns true for files that should NOT be extracted to data/plugins/.
func isSkippedArchivePath(relPath string) bool {
	parts := strings.Split(relPath, "/")
	for _, p := range parts {
		switch p {
		case ".git", ".github", "node_modules", ".idea", ".vscode":
			return true
		}
	}
	base := path.Base(relPath)
	switch base {
	case ".gitignore", ".gitattributes", "go.mod", "go.sum":
		return true
	}
	return false
}

// detectArchivePrefix finds the common directory prefix in archive paths.
// GitHub zipball wraps all files under "{owner}-{repo}-{sha}/".
func detectArchivePrefix(files map[string][]byte) string {
	var prefix string
	for name := range files {
		idx := strings.IndexByte(name, '/')
		if idx < 0 {
			return "" // file at root — no prefix
		}
		candidate := name[:idx+1]
		if prefix == "" {
			prefix = candidate
		} else if candidate != prefix {
			return "" // inconsistent — no common prefix
		}
	}
	return prefix
}

// pluginDependency is the minimal dependency spec carried from plugin.yaml
// through the install pipeline. Mirrors sdk.Dependency but avoids importing
// the SDK package in the logic layer.
type pluginDependency struct {
	ID       string
	Version  string
	Optional bool
}

// parseArchiveResult holds all data extracted from a plugin archive.
type parseArchiveResult struct {
	Manifest  *eng.Manifest
	Script    string
	Assets    map[string][]byte
	AllFiles  map[string][]byte  // all files in archive (for JS/builtin plugins that need full extraction)
	Homepage  string             // from package.json "homepage" field, used as repo_url fallback
	IsJS      bool               // true when plugin type is "js" or "full"
	IsBuiltin bool               // true when plugin type is "builtin" (compiled Go)
	JSEntry   string             // JS entry file name from plugin.yaml (default "plugin.js")
	Depends   []pluginDependency // from plugin.yaml depends: list (empty for package.json manifests)
}

// pluginYAMLManifest is the YAML structure for plugin.yaml files.
type pluginYAMLManifest struct {
	ID          string `yaml:"id"`
	Title       string `yaml:"title"`
	Version     string `yaml:"version"`
	Icon        string `yaml:"icon"`
	Author      string `yaml:"author"`
	Description string `yaml:"description"`
	License     string `yaml:"license"`
	TrustLevel  string `yaml:"trust_level"`
	SDKVersion  string `yaml:"sdk_version"`
	Homepage    string `yaml:"homepage"`

	Type    string `yaml:"type"`     // builtin, js, yaml, ui, full
	Runtime string `yaml:"runtime"`  // compiled | interpreted
	Bundled bool   `yaml:"bundled"`  // included in official prebuilt binary
	JSEntry string `yaml:"js_entry"` // JS entry file (default "plugin.js")
	Layer   any    `yaml:"layer"`    // int or []int
	Entry   string `yaml:"entry"`
	// (admin_js/public_js/css removed in manifest v2 — now in contributes)

	Settings []struct {
		Key         string   `yaml:"key"`
		Label       string   `yaml:"label"`
		Type        string   `yaml:"type"`
		Required    bool     `yaml:"required"`
		Default     any      `yaml:"default"`
		Placeholder string   `yaml:"placeholder"`
		Description string   `yaml:"description"`
		Options     []string `yaml:"options"`
		Group       string   `yaml:"group"`
		Shared      bool     `yaml:"shared"`
	} `yaml:"settings"`

	Pages []struct {
		Path      string `yaml:"path"`
		Slot      string `yaml:"slot"`
		Component string `yaml:"component"`
		Title     string `yaml:"title"`
		Nav       *struct {
			Type  string `yaml:"type"`
			Group string `yaml:"group"`
			Icon  string `yaml:"icon"`
			Order int    `yaml:"order"`
		} `yaml:"nav"`
	} `yaml:"pages"`

	Contributes *struct {
		Commands []struct {
			ID    string `yaml:"id"`
			Title string `yaml:"title"`
			Icon  string `yaml:"icon"`
		} `yaml:"commands"`
		Menus map[string][]struct {
			Command string `yaml:"command"`
		} `yaml:"menus"`
	} `yaml:"contributes"`

	I18n map[string]map[string]string `yaml:"i18n"`

	Migrations []struct {
		Version int               `yaml:"version"`
		Up      map[string]string `yaml:"up"`
		Down    map[string]string `yaml:"down"`
	} `yaml:"migrations"`

	Routes []struct {
		Method      string `yaml:"method"`
		Path        string `yaml:"path"`
		Auth        string `yaml:"auth"`
		Fn          string `yaml:"fn"`
		Description string `yaml:"description"`
	} `yaml:"routes"`

	Depends []struct {
		ID       string `yaml:"id"`
		Version  string `yaml:"version"`
		Optional bool   `yaml:"optional"`
	} `yaml:"depends"`
}

// parseArchive reads the plugin manifest and entry script from any archive format
// supported by mholt/archives (zip, tar.gz, tar.bz2, tar.xz, 7z, rar, …).
//
// Manifest resolution order:
//  1. plugin.yaml — new declarative format (preferred)
//  2. package.json with a "plugin" field — legacy npm-style manifest
func parseArchive(data []byte) (*parseArchiveResult, error) {
	ctx := context.Background()

	format, reader, err := archives.Identify(ctx, "", bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("unsupported archive format: %w", err)
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return nil, fmt.Errorf("archive format %q cannot be extracted", format.Extension())
	}

	// Collect every file in the archive into a map (clean path → content).
	collected := make(map[string][]byte)
	if err = ex.Extract(ctx, reader, func(ctx context.Context, f archives.FileInfo) error {
		if f.IsDir() {
			return nil
		}
		rc, e := f.Open()
		if e != nil {
			return fmt.Errorf("open %s: %w", f.NameInArchive, e)
		}
		defer rc.Close()
		b, e := io.ReadAll(rc)
		if e != nil {
			return fmt.Errorf("read %s: %w", f.NameInArchive, e)
		}
		collected[path.Clean(f.NameInArchive)] = b
		return nil
	}); err != nil {
		return nil, fmt.Errorf("extract archive: %w", err)
	}

	// ── Strategy 1: Try plugin.yaml first ────────────────────────────────────

	var m eng.Manifest
	var homepage string
	var yamlJSEntry string
	var yamlType string
	var yamlDepends []pluginDependency
	resolved := false

	for name, content := range collected {
		base := path.Base(name)
		if base != "plugin.yaml" && base != "plugin.yml" {
			continue
		}
		// Parse raw YAML to extract js_entry, type, and depends before converting to Manifest
		var raw pluginYAMLManifest
		if yaml.Unmarshal(content, &raw) == nil {
			yamlJSEntry = raw.JSEntry
			yamlType = raw.Type
			for _, d := range raw.Depends {
				yamlDepends = append(yamlDepends, pluginDependency{
					ID: d.ID, Version: d.Version, Optional: d.Optional,
				})
			}
		}
		if mf, hp, err := parsePluginYAML(content); err == nil {
			m = *mf
			homepage = hp
			resolved = true
			break
		}
	}

	// ── Strategy 2: Fallback to package.json "plugin" field ──────────────────

	if !resolved {
		for name, content := range collected {
			if path.Base(name) != "package.json" {
				continue
			}
			var pkg struct {
				Name        string `json:"name"`
				Version     string `json:"version"`
				Description string `json:"description"`
				Author      any    `json:"author"`
				Homepage    string `json:"homepage"`
				Plugin      *struct {
					Title            string                 `json:"title"`
					Icon             string                 `json:"icon"`
					Entry            string                 `json:"entry"`
					CSS              string                 `json:"css"`
					Priority         int                    `json:"priority"`
					Settings         []eng.SettingField     `json:"settings"`
					Capabilities     eng.PluginCapabilities `json:"capabilities"`
					Webhooks         []eng.WebhookDef       `json:"webhooks"`
					Pipelines        []eng.PipelineDef      `json:"pipelines"`
					MinHostVersion   string                 `json:"minHostVersion"`
					TrustLevel       eng.TrustLevel         `json:"trust_level"`
					ActivationEvents []string               `json:"activationEvents"`
					// (admin_js/public_js removed in manifest v2)
					Routes           []eng.RouteDef         `json:"routes"`
					Contributes      *eng.Contributes       `json:"contributes"`
					Migrations       []eng.MigrationDef     `json:"migrations"`
					Pages            []eng.PageDef          `json:"pages"`
					Service          *eng.ServiceDef        `json:"service"`
				} `json:"plugin"`
			}
			if err = json.Unmarshal(content, &pkg); err != nil || pkg.Plugin == nil {
				continue
			}
			m.Name = pkg.Name
			m.Version = pkg.Version
			m.Description = pkg.Description
			switch v := pkg.Author.(type) {
			case string:
				m.Author = v
			case map[string]any:
				if s, ok := v["name"].(string); ok {
					m.Author = s
				}
			}
			m.Title = pkg.Plugin.Title
			m.Icon = pkg.Plugin.Icon
			m.Entry = pkg.Plugin.Entry
			m.Priority = pkg.Plugin.Priority
			m.Settings = pkg.Plugin.Settings
			m.Capabilities = pkg.Plugin.Capabilities
			m.Webhooks = pkg.Plugin.Webhooks
			m.Pipelines = pkg.Plugin.Pipelines
			m.MinHostVersion = pkg.Plugin.MinHostVersion
			m.TrustLevel = pkg.Plugin.TrustLevel
			m.ActivationEvents = pkg.Plugin.ActivationEvents
			// (admin_js/public_js now in contributes)
			m.Routes = pkg.Plugin.Routes
			m.Contributes = pkg.Plugin.Contributes
			m.Migrations = pkg.Plugin.Migrations
			m.Service = pkg.Plugin.Service
			homepage = pkg.Homepage
			resolved = true
			break
		}
	}

	if !resolved || m.Name == "" {
		return nil, fmt.Errorf("manifest not found: archive must contain plugin.yaml or package.json with a \"plugin\" field")
	}

	// Determine if this is a JS (Goja) plugin
	isJS := yamlType == "js" || yamlType == "full"
	if !isJS {
		// Also detect by presence of plugin.js in archive
		for name := range collected {
			if path.Base(name) == "plugin.js" {
				isJS = true
				break
			}
		}
	}

	// For JS/full plugins, preserve all files for extraction to disk
	if isJS {
		return &parseArchiveResult{
			Manifest: &m,
			AllFiles: collected,
			Homepage: homepage,
			IsJS:     true,
			JSEntry:  yamlJSEntry,
			Depends:  yamlDepends,
		}, nil
	}

	// Builtin (compiled Go) plugins: preserve all files so frontend assets
	// can be extracted to data/plugins/{id}/. The Go backend logic is compiled
	// into the binary and will be activated on next server restart.
	if yamlType == "builtin" {
		return &parseArchiveResult{
			Manifest:  &m,
			AllFiles:  collected,
			Homepage:  homepage,
			IsBuiltin: true,
			Depends:   yamlDepends,
		}, nil
	}

	// ── JS/YAML plugin path: locate entry script ─────────────────────────────
	if m.Entry == "" {
		m.Entry = "dist/index.js"
	}

	entryBase := path.Base(m.Entry)
	var script string
	for name, content := range collected {
		if name == m.Entry || path.Base(name) == entryBase {
			script = string(content)
			break
		}
	}
	// Layer 2/3 plugins may have no Goja entry script — that's OK if they have frontend assets
	hasActivation := m.Contributes != nil && len(m.Contributes.Activation) > 0
	hasPages := m.Contributes != nil && len(m.Contributes.Pages) > 0
	if script == "" && !hasActivation && !hasPages {
		return nil, fmt.Errorf("entry file %s not found in archive", m.Entry)
	}

	// ── Collect frontend asset files ─────────────────────────────────────────
	assets := make(map[string][]byte)
	allowedAssetExts := map[string]bool{".js": true, ".mjs": true, ".css": true, ".json": true}
	for name, content := range collected {
		ext := path.Ext(name)
		if !allowedAssetExts[ext] {
			continue
		}
		base := path.Base(name)
		if base == entryBase || base == "package.json" || base == "plugin.yaml" || base == "plugin.yml" {
			continue
		}
		assets[base] = content
	}

	return &parseArchiveResult{
		Manifest: &m,
		Script:   script,
		Assets:   assets,
		Homepage: homepage,
		Depends:  yamlDepends,
	}, nil
}

// parsePluginYAML converts a plugin.yaml content into an engine Manifest.
func parsePluginYAML(data []byte) (*eng.Manifest, string, error) {
	var y pluginYAMLManifest
	if err := yaml.Unmarshal(data, &y); err != nil {
		return nil, "", err
	}
	if y.ID == "" {
		return nil, "", fmt.Errorf("plugin.yaml missing required 'id' field")
	}

	m := &eng.Manifest{
		Name:        y.ID,
		Title:       y.Title,
		Version:     y.Version,
		Icon:        y.Icon,
		Author:      y.Author,
		Description: y.Description,
		TrustLevel:  eng.TrustLevel(y.TrustLevel),
		Type:        y.Type,
		Entry:       y.Entry,
		I18n:        y.I18n,
	}

	// Settings
	for _, s := range y.Settings {
		sf := eng.SettingField{
			Key:         s.Key,
			Label:       s.Label,
			Type:        eng.SettingType(s.Type),
			Required:    s.Required,
			Default:     s.Default,
			Placeholder: s.Placeholder,
			Description: s.Description,
			Options:     s.Options,
			Group:       s.Group,
			Shared:      s.Shared,
		}
		m.Settings = append(m.Settings, sf)
	}

	// Contributes (pages are now inside contributes)
	c := m.Contributes
	if c == nil && (len(y.Pages) > 0 || y.Contributes != nil) {
		c = &eng.Contributes{}
	}
	if y.Contributes != nil {
		for _, cmd := range y.Contributes.Commands {
			c.Commands = append(c.Commands, eng.CommandDef{
				ID:    cmd.ID,
				Title: cmd.Title,
				Icon:  cmd.Icon,
			})
		}
		if len(y.Contributes.Menus) > 0 {
			c.Menus = make(map[string][]eng.MenuEntry)
			for slot, entries := range y.Contributes.Menus {
				for _, e := range entries {
					c.Menus[slot] = append(c.Menus[slot], eng.MenuEntry{Command: e.Command})
				}
			}
		}
	}
	// Pages
	for _, p := range y.Pages {
		pd := eng.PageDef{
			Path:      p.Path,
			Slot:      p.Slot,
			Component: p.Component,
			Title:     p.Title,
		}
		if p.Nav != nil {
			pd.Nav = &eng.NavDef{
				Type:  p.Nav.Type,
				Group: p.Nav.Group,
				Icon:  p.Nav.Icon,
				Order: p.Nav.Order,
			}
		}
		c.Pages = append(c.Pages, pd)
	}
	m.Contributes = c

	// Migrations
	for _, mig := range y.Migrations {
		m.Migrations = append(m.Migrations, eng.MigrationDef{
			Version: mig.Version,
			Up:      mig.Up,
			Down:    mig.Down,
		})
	}

	// Routes
	for _, r := range y.Routes {
		m.Routes = append(m.Routes, eng.RouteDef{
			Method:      r.Method,
			Path:        r.Path,
			Auth:        r.Auth,
			Fn:          r.Fn,
			Description: r.Description,
		})
	}

	return m, y.Homepage, nil
}

// ── Plugin Assets ───────────────────────────────────────────────────────────

// PluginAssetsDir returns the directory where plugin frontend assets are stored.
func PluginAssetsDir() string { return filepath.Join("data", "plugins") }

// savePluginAssets writes frontend asset files to data/plugins/{sanitized-id}/.
func savePluginAssets(pluginID string, assets map[string][]byte) error {
	dir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(pluginID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create assets dir: %w", err)
	}
	for name, content := range assets {
		fp := filepath.Join(dir, name)
		if err := os.WriteFile(fp, content, 0644); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}
	return nil
}

// sanitizePluginPath converts a plugin ID like "nuxtblog/ai-polish" into
// a safe directory name "nuxtblog--ai-polish" to avoid nested directories.
func sanitizePluginPath(id string) string {
	return strings.ReplaceAll(id, "/", "--")
}
