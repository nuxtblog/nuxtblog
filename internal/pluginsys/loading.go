package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/internal/logic/payment"
	"github.com/nuxtblog/nuxtblog/sdk"
)

// activatePlugin runs the full activation lifecycle for a sdk.Plugin instance.
func (m *Manager) activatePlugin(ctx context.Context, p sdk.Plugin, source string, pluginDir, jsFile string) error {
	mf := p.Manifest()
	id := mf.ID

	lp := &loadedPlugin{
		plugin:    p,
		pluginDir: pluginDir,
		jsFile:    jsFile,
		stats:     &PluginStats{},
		window:    &SlidingWindow{},
		errors:    &ErrorRingBuffer{},
	}

	// Build plugin context with DB capabilities
	var dbCaps *DBCap
	if mf.Capabilities != nil && mf.Capabilities.DB != nil {
		dbCaps = convertSDKDBCap(mf.Capabilities.DB)
	}
	trust := TrustLevel(mf.TrustLevel)
	if source == "builtin" {
		trust = TrustLevelOfficial
	}

	pctx := sdk.PluginContext{
		DB:       newPluginDB(id, dbCaps, trust),
		Store:    newPluginStore(id),
		Settings: newPluginSettings(id),
		Log:      newPluginLogger(id),
		Plugins:  &pluginQuery{mgr: m},
		AI:       DefaultAI(),
		Media:    newMediaService(id),
		Commerce: newPluginCommerce(id),
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

	// Register media categories declared in capabilities.media
	if mf.Capabilities != nil && mf.Capabilities.Media != nil {
		for _, cat := range mf.Capabilities.Media.Categories {
			if err := pctx.Media.RegisterCategory(cat); err != nil {
				g.Log().Warningf(ctx, "[pluginmgr] %s: register media category %s: %v", id, cat.Slug, err)
			}
		}
	}

	// Register payment gateway if plugin provides one
	if gw, ok := p.(sdk.HasPaymentGateway); ok {
		payment.RegisterPluginGateway(gw.PaymentGateway())
		g.Log().Infof(ctx, "[pluginmgr] plugin %s: registered payment gateway %s", id, gw.PaymentGateway().ProviderInfo().Slug)
	}

	// Register provider interfaces
	if wp, ok := p.(sdk.HasWalletProvider); ok {
		m.walletProvider = wp.WalletProvider()
		g.Log().Infof(ctx, "[pluginmgr] plugin %s: registered wallet provider", id)
	}
	if cp, ok := p.(sdk.HasCreditsProvider); ok {
		m.creditsProvider = cp.CreditsProvider()
		g.Log().Infof(ctx, "[pluginmgr] plugin %s: registered credits provider", id)
	}
	if mp, ok := p.(sdk.HasMembershipProvider); ok {
		m.membershipProvider = mp.MembershipProvider()
		g.Log().Infof(ctx, "[pluginmgr] plugin %s: registered membership provider", id)
	}
	if ep, ok := p.(sdk.HasEntitlementProvider); ok {
		m.entitlementProvider = ep.EntitlementProvider()
		g.Log().Infof(ctx, "[pluginmgr] plugin %s: registered entitlement provider", id)
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

	// Check dependencies before loading
	for _, dep := range pm.Depends {
		if !dep.Optional && !m.HasPlugin(dep.ID) {
			return fmt.Errorf("missing required dependency: %s", dep.ID)
		}
		if dep.Version != "" && m.HasPlugin(dep.ID) {
			ver := (&pluginQuery{mgr: m}).GetVersion(dep.ID)
			if ver != "" && !matchSemverConstraint(ver, dep.Version) {
				return fmt.Errorf("dependency %s requires version %s, but found %s", dep.ID, dep.Version, ver)
			}
		}
	}

	// Capture source info of dependents that will be cascade-unloaded, so they
	// can be reloaded afterwards. Order: shallowest first (reverse of cascade
	// unload order), excluding pm.ID itself.
	type reloadTarget struct {
		id        string
		pluginDir string
		jsFile    string
	}
	var toReload []reloadTarget
	m.mu.RLock()
	cascade := m.graph.GetCascadeUnloadOrder(pm.ID) // [deepest…, pm.ID]
	for i := len(cascade) - 1; i >= 0; i-- {
		depID := cascade[i]
		if depID == pm.ID {
			continue
		}
		if lp, ok := m.plugins[depID]; ok && lp.pluginDir != "" {
			toReload = append(toReload, reloadTarget{
				id:        depID,
				pluginDir: lp.pluginDir,
				jsFile:    lp.jsFile,
			})
		}
	}
	m.mu.RUnlock()

	// Unload previous version if already loaded (update scenario)
	m.UnloadPlugin(pm.ID)

	m.graph.Add(pm.ID, pm.Depends)

	p, err := loadGojaPlugin(ctx, pluginDir, jsFile, *pm)
	if err != nil {
		return fmt.Errorf("goja load: %w", err)
	}
	if err := m.activatePlugin(ctx, p, "external", pluginDir, jsFile); err != nil {
		return err
	}
	m.ensureDBRecordExt(ctx, *pm, "external")

	// Cascade reload dependents in topological order (shallow → deep). Each
	// recursive call's own cascade capture will be empty because deeper
	// dependents are still unloaded at that point — no double-reload.
	for _, rt := range toReload {
		if err := m.InstallJSPlugin(ctx, rt.pluginDir, rt.jsFile); err != nil {
			g.Log().Warningf(ctx,
				"[pluginmgr] cascade reload dependent %s failed: %v", rt.id, err)
			continue
		}
		g.Log().Infof(ctx, "[pluginmgr] cascade reloaded dependent: %s", rt.id)
	}

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
			"capabilities":    manifestCapsJSON(mf),
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

// ensureDBRecordDiscovered inserts a discovered-but-not-yet-loaded plugin into
// the DB with enabled=0, so it appears in the admin panel for manual activation.
// Does nothing if the plugin already exists in the DB.
func (m *Manager) ensureDBRecordDiscovered(ctx context.Context, mf sdk.Manifest) {
	db := g.DB()
	count, _ := db.Ctx(ctx).Model("plugins").Where("id", mf.ID).Count()
	if count > 0 {
		return // already in DB, respect existing enabled state
	}

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

	_, _ = db.Ctx(ctx).Model("plugins").Data(g.Map{
		"id":              mf.ID,
		"title":           mf.Title,
		"description":     mf.Description,
		"version":         mf.Version,
		"author":          mf.Author,
		"icon":            mf.Icon,
		"enabled":         0, // disabled by default — admin must enable manually
		"source":          "external",
		"script":          "",
		"settings":        defaultsJSON,
		"settings_schema": schemaJSON,
		"manifest":        manifestJSON,
		"capabilities":    manifestCapsJSON(mf),
		"installed_at":    time.Now(),
	}).Insert()

	g.Log().Infof(ctx, "[pluginmgr] discovered new plugin %s, registered as disabled", mf.ID)
}

// registerMetadataPlugin registers a YAML/UI-only plugin in the DB (no runtime).
func (m *Manager) registerMetadataPlugin(ctx context.Context, mf *sdk.Manifest, source string) {
	m.ensureDBRecordExt(ctx, *mf, source)
	g.Log().Infof(ctx, "[pluginmgr] registered %s plugin: %s v%s", mf.Type, mf.ID, mf.Version)
}
