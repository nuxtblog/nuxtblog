// Package pluginsys — loader.go handles plugin discovery, activation, and DB persistence.
//
// External plugins are loaded from data/plugins/ directories.
// JS plugins (type: js/full) are interpreted by Goja.
// Builtin plugins are compiled Go, loaded via LoadStatic in manager.go.
package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// LoadExternal scans dataDir for plugin directories, reads plugin.yaml,
// and loads JS/full type plugins via Goja.
//
// dataDir is typically "data/plugins/".
func (m *Manager) LoadExternal(ctx context.Context, dataDir string) error {
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // no external plugins directory
		}
		return fmt.Errorf("read plugins dir: %w", err)
	}

	g.Log().Infof(ctx, "[pluginmgr] LoadExternal scanning %s (%d entries)", dataDir, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pluginDir := filepath.Join(dataDir, entry.Name())

		// Read plugin.yaml
		yamlPath := filepath.Join(pluginDir, "plugin.yaml")
		if _, err := os.Stat(yamlPath); err != nil {
			yamlPath = filepath.Join(pluginDir, "plugin.yml")
			if _, err := os.Stat(yamlPath); err != nil {
				continue // no manifest, skip
			}
		}

		yamlData, err := os.ReadFile(yamlPath)
		if err != nil {
			g.Log().Warningf(ctx, "[pluginmgr] cannot read %s: %v", yamlPath, err)
			continue
		}

		pm, err := parsePluginYAML(yamlData)
		if err != nil {
			g.Log().Warningf(ctx, "[pluginmgr] bad plugin.yaml in %s: %v", pluginDir, err)
			continue
		}

		g.Log().Infof(ctx, "[pluginmgr] found external plugin: %s (type=%s)", pm.ID, pm.Type)

		// Skip if already loaded (builtin takes priority)
		m.mu.RLock()
		_, loaded := m.plugins[pm.ID]
		m.mu.RUnlock()
		if loaded {
			g.Log().Infof(ctx, "[pluginmgr] %s already loaded (builtin), skipping", pm.ID)
			continue
		}

		// Check if enabled
		if !m.isEnabled(ctx, pm.ID) {
			g.Log().Infof(ctx, "[pluginmgr] external plugin %s not enabled, skipping", pm.ID)
			continue
		}

		switch pm.Type {
		case sdk.TypeJS, sdk.TypeFull:
			jsFile := pm.JSEntry
			if jsFile == "" {
				jsFile = "plugin.js"
			}
			jsPath := filepath.Join(pluginDir, jsFile)
			if _, err := os.Stat(jsPath); err != nil {
				g.Log().Warningf(ctx, "[pluginmgr] %s: js source %s not found", pm.ID, jsFile)
				continue
			}
			p, err := loadGojaPlugin(ctx, pluginDir, jsFile, *pm)
			if err != nil {
				g.Log().Errorf(ctx, "[pluginmgr] %s goja load failed: %v", pm.ID, err)
				continue
			}
			if err := m.activatePlugin(ctx, p, "external"); err != nil {
				g.Log().Errorf(ctx, "[pluginmgr] %s activation failed: %v", pm.ID, err)
			}
			m.ensureDBRecordExt(ctx, *pm, "external")

		case sdk.TypeBuiltin:
			// Builtin plugins have Go backend compiled into the binary.
			// Here we only register metadata (frontend assets, pages, contributes)
			// so the admin panel can serve the plugin's UI components.
			m.registerMetadataPlugin(ctx, pm, "builtin")

		case sdk.TypeYAML:
			m.registerMetadataPlugin(ctx, pm, "external")

		case sdk.TypeUI:
			m.registerMetadataPlugin(ctx, pm, "external")

		default:
			g.Log().Infof(ctx, "[pluginmgr] %s: type '%s' not handled by Go loader", pm.ID, pm.Type)
		}
	}
	return nil
}

// activatePlugin runs the full activation lifecycle for a sdk.Plugin instance.
func (m *Manager) activatePlugin(ctx context.Context, p sdk.Plugin, source string) error {
	mf := p.Manifest()
	id := mf.ID

	lp := &loadedPlugin{
		plugin: p,
		stats:  &PluginStats{},
		window: &SlidingWindow{},
		errors: &ErrorRingBuffer{},
	}

	// Build plugin context
	pctx := sdk.PluginContext{
		DB:       newPluginDB(id),
		Store:    newPluginStore(id),
		Settings: newPluginSettings(id),
		Log:      newPluginLogger(id),
	}
	lp.ctx = pctx

	// Run migrations
	if migs := p.Migrations(); len(migs) > 0 {
		if err := m.runMigrations(ctx, id, migs); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}

	// Activate
	if err := p.Activate(pctx); err != nil {
		return fmt.Errorf("activate: %w", err)
	}

	// Collect routes
	reg := &registrar{pluginID: id}
	p.Routes(reg)
	lp.routes = reg.entries

	// Collect filters
	lp.filters = p.Filters()

	m.mu.Lock()
	m.plugins[id] = lp
	m.mu.Unlock()

	// Dynamically register routes on the running HTTP server (hot install)
	if m.server != nil && len(lp.routes) > 0 {
		m.bindPluginRoutes(ctx, m.server, id, lp)
	}

	// Write embedded frontend assets to data/plugins/{id}/ for serving.
	// Builtin plugins skip this: their assets are //go:embed'd into the binary
	// and served directly from memory by RegisterAssetRoutes (assets.go).
	if source != "builtin" {
		if ap, ok := p.(sdk.HasAssets); ok {
			if assets := ap.Assets(); len(assets) > 0 {
				assetsDir := filepath.Join("data", "plugins", id)
				if err := os.MkdirAll(assetsDir, 0o755); err != nil {
					g.Log().Warningf(ctx, "[pluginmgr] plugin %s: mkdir assets: %v", id, err)
				} else {
					for name, data := range assets {
						fp := filepath.Join(assetsDir, filepath.Base(name))
						if err := os.WriteFile(fp, data, 0o644); err != nil {
							g.Log().Warningf(ctx, "[pluginmgr] plugin %s: write asset %s: %v", id, name, err)
						}
					}
					g.Log().Infof(ctx, "[pluginmgr] plugin %s: wrote %d frontend assets", id, len(assets))
				}
			}
		}
	}

	// Register in DB
	m.ensureDBRecordExt(ctx, mf, source)

	g.Log().Infof(ctx, "[pluginmgr] loaded external plugin: %s v%s", id, mf.Version)
	return nil
}

// ensureDBRecordExt inserts or updates a plugin's metadata in the DB.
func (m *Manager) ensureDBRecordExt(ctx context.Context, mf sdk.Manifest, source string) {
	db := g.DB()
	id := mf.ID

	schemaJSON := "[]"
	if len(mf.Settings) > 0 {
		if b, err := json.Marshal(mf.Settings); err == nil {
			schemaJSON = string(b)
		}
	}

	defaults := make(map[string]any, len(mf.Settings))
	for _, s := range mf.Settings {
		if s.Default != nil {
			defaults[s.Key] = s.Default
		}
	}
	defaultsJSON := "{}"
	if b, err := json.Marshal(defaults); err == nil {
		defaultsJSON = string(b)
	}

	manifestJSON := "{}"
	if b, err := json.Marshal(mf); err == nil {
		manifestJSON = string(b)
	}

	count, _ := db.Ctx(ctx).Model("plugins").Where("id", id).Count()
	if count == 0 {
		_, _ = db.Ctx(ctx).Model("plugins").Data(g.Map{
			"id":              id,
			"title":           mf.Title,
			"description":     mf.Description,
			"version":         mf.Version,
			"author":          mf.Author,
			"icon":            mf.Icon,
			"enabled":         1,
			"source":          source,
			"script":          "",
			"settings":        defaultsJSON,
			"settings_schema": schemaJSON,
			"manifest":        manifestJSON,
			"capabilities":    `{"store":{"read":true,"write":true},"db":true}`,
			"installed_at":    time.Now(),
		}).Insert()
	} else {
		_, _ = db.Ctx(ctx).Model("plugins").Where("id", id).Data(g.Map{
			"title":           mf.Title,
			"description":     mf.Description,
			"version":         mf.Version,
			"author":          mf.Author,
			"icon":            mf.Icon,
			"source":          source,
			"settings_schema": schemaJSON,
			"manifest":        manifestJSON,
		}).Update()
	}
}

// registerMetadataPlugin registers a YAML/UI-only plugin in the DB (no runtime).
func (m *Manager) registerMetadataPlugin(ctx context.Context, mf *sdk.Manifest, source string) {
	m.ensureDBRecordExt(ctx, *mf, source)
	g.Log().Infof(ctx, "[pluginmgr] registered %s plugin: %s v%s", mf.Type, mf.ID, mf.Version)
}

// InstallJSPlugin loads a Goja JS plugin from a directory and activates it.
// Called by the plugin installer after extracting files to data/plugins/{id}/.
func (m *Manager) InstallJSPlugin(ctx context.Context, pluginDir string, jsFile string) error {
	if jsFile == "" {
		jsFile = "plugin.js"
	}
	// Read manifest from plugin.yaml
	yamlPath := filepath.Join(pluginDir, "plugin.yaml")
	yamlData, err := os.ReadFile(yamlPath)
	if err != nil {
		return fmt.Errorf("read plugin.yaml: %w", err)
	}
	pm, err := parsePluginYAML(yamlData)
	if err != nil {
		return fmt.Errorf("parse plugin.yaml: %w", err)
	}

	// Unload previous version if already loaded (update scenario)
	m.UnloadPlugin(pm.ID)

	p, err := loadGojaPlugin(ctx, pluginDir, jsFile, *pm)
	if err != nil {
		return fmt.Errorf("goja load: %w", err)
	}
	if err := m.activatePlugin(ctx, p, "external"); err != nil {
		return err
	}
	m.ensureDBRecordExt(ctx, *pm, "external")
	return nil
}

// UnloadPlugin deactivates and removes a plugin from the manager.
func (m *Manager) UnloadPlugin(id string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if lp, ok := m.plugins[id]; ok {
		if dp, ok := lp.plugin.(sdk.HasDeactivate); ok {
			_ = dp.Deactivate()
		}
		delete(m.plugins, id)
	}
}
