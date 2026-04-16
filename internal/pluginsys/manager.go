// Package pluginsys manages Go-native plugin lifecycle.
//
// It loads statically compiled plugins (via sdk.GetStatic()) and
// dynamically interpreted plugins (via Goja JS) then wires them
// into the platform's route system, event bus, and filter chain.
package pluginsys

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// Manager holds all loaded Go-native plugins.
type Manager struct {
	mu      sync.RWMutex
	plugins map[string]*loadedPlugin
	server  *ghttp.Server // set by RegisterRoutes, used for dynamic route binding
	graph   *depGraph
	eventWg sync.WaitGroup // tracks in-flight FanOutEvent goroutines

	// Provider references — set during plugin activation
	walletProvider      sdk.WalletProvider
	creditsProvider     sdk.CreditsProvider
	membershipProvider  sdk.MembershipProvider
	entitlementProvider sdk.EntitlementProvider
}

// WalletProvider returns the registered wallet provider (may be nil).
func (m *Manager) WalletProvider() sdk.WalletProvider { return m.walletProvider }

// CreditsProvider returns the registered credits provider (may be nil).
func (m *Manager) CreditsProvider() sdk.CreditsProvider { return m.creditsProvider }

// MembershipProvider returns the registered membership provider (may be nil).
func (m *Manager) MembershipProvider() sdk.MembershipProvider { return m.membershipProvider }

// EntitlementProvider returns the registered entitlement provider (may be nil).
func (m *Manager) EntitlementProvider() sdk.EntitlementProvider { return m.entitlementProvider }

type loadedPlugin struct {
	plugin  sdk.Plugin
	ctx     sdk.PluginContext
	routes  []routeEntry
	filters []sdk.FilterDef

	// Source info — needed for cascade reload after a dependency is updated.
	// Empty for builtin plugins (they don't need on-disk reloading).
	pluginDir string
	jsFile    string

	// Observability
	stats  *PluginStats
	window *SlidingWindow
	errors *ErrorRingBuffer
}

type routeEntry struct {
	method  string
	path    string
	handler http.HandlerFunc
	auth    string
}

// defaultMgr is the global plugin manager, set by New().
var defaultMgr *Manager

// New creates a new plugin manager and stores it as the global default.
func New() *Manager {
	m := &Manager{
		plugins: make(map[string]*loadedPlugin),
		graph:   newDepGraph(),
	}
	defaultMgr = m
	return m
}

// GetManager returns the global plugin manager instance.
func GetManager() *Manager {
	return defaultMgr
}

// HasPlugin reports whether a Go-native plugin with the given ID is registered.
// Checks both already-loaded plugins and statically registered (not yet loaded) ones.
func (m *Manager) HasPlugin(id string) bool {
	// Check loaded plugins
	m.mu.RLock()
	_, ok := m.plugins[id]
	m.mu.RUnlock()
	if ok {
		return true
	}
	// Check static registry (plugins registered via init() but not yet loaded)
	for _, p := range sdk.GetStatic() {
		if p.Manifest().ID == id {
			return true
		}
	}
	return false
}

// LoadStatic loads all statically registered plugins and activates them.
// Uses two-phase loading: first collects all enabled candidates, registers
// dependencies, topologically sorts, then activates in dependency order.
func (m *Manager) LoadStatic(ctx context.Context) error {
	// Phase 1: sync DB records for ALL builtin plugins (even disabled ones)
	// so their manifest/contributes stay up-to-date in the admin UI.
	// Then collect enabled candidates for activation.
	candidates := make(map[string]sdk.Plugin)
	var ids []string
	for _, p := range sdk.GetStatic() {
		mf := p.Manifest()
		id := mf.ID
		if id == "" {
			g.Log().Warning(ctx, "[pluginmgr] skipping plugin with empty ID")
			continue
		}
		m.ensureDBRecord(ctx, mf)
		if !m.isEnabled(ctx, id) {
			g.Log().Infof(ctx, "[pluginmgr] plugin %s not enabled, skipping", id)
			continue
		}
		candidates[id] = p
		ids = append(ids, id)
		m.graph.Add(id, mf.Depends)
	}

	// Phase 2: topological sort
	versionResolver := func(id string) string {
		if p, ok := candidates[id]; ok {
			return p.Manifest().Version
		}
		return ""
	}
	sorted, err := m.graph.TopologicalSort(ids, versionResolver)
	if err != nil {
		g.Log().Warningf(ctx, "[pluginmgr] dependency sort failed: %v; loading in original order", err)
		sorted = ids
	}

	// Phase 3: activate in sorted order
	for _, id := range sorted {
		p := candidates[id]
		if p == nil {
			continue
		}
		if err := m.activatePlugin(ctx, p, "builtin", "", ""); err != nil {
			g.Log().Errorf(ctx, "[pluginmgr] plugin %s failed: %v", id, err)
			continue
		}
	}
	return nil
}

// RunFilter runs all Go-native plugin filters for the given event.
func (m *Manager) RunFilter(ctx context.Context, event string, data, meta map[string]any) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, lp := range m.plugins {
		for _, f := range lp.filters {
			if f.Event != event {
				continue
			}
			fc := &sdk.FilterContext{
				Context: ctx,
				Event:   event,
				Data:    data,
				Meta:    meta,
			}
			start := time.Now()
			f.Handler(fc)
			dur := time.Since(start)

			var filterErr error
			if fc.IsAborted() {
				filterErr = fmt.Errorf("filter aborted by plugin: %s", fc.AbortReason())
			}
			lp.recordExec(id, "filter:"+event, dur, filterErr)

			if filterErr != nil {
				return filterErr
			}
		}
	}
	return nil
}

// FanOutEvent dispatches an event to all Go-native plugins that implement HasEvents.
func (m *Manager) FanOutEvent(ctx context.Context, event string, data map[string]any) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, lp := range m.plugins {
		if ep, ok := lp.plugin.(sdk.HasEvents); ok {
			m.eventWg.Add(1)
			go func(id string, lp *loadedPlugin, ep sdk.HasEvents) {
				defer m.eventWg.Done()
				start := time.Now()
				ep.OnEvent(ctx, event, data)
				dur := time.Since(start)
				lp.recordExec(id, "event:"+event, dur, nil)
			}(id, lp, ep)
		}
	}
}

// Shutdown waits for in-flight event handlers, then deactivates all loaded plugins.
func (m *Manager) Shutdown() {
	m.eventWg.Wait()

	m.mu.Lock()
	defer m.mu.Unlock()

	for id, lp := range m.plugins {
		if dp, ok := lp.plugin.(sdk.HasDeactivate); ok {
			if err := dp.Deactivate(); err != nil {
				g.Log().Warningf(context.Background(), "[pluginmgr] plugin %s deactivate error: %v", id, err)
			}
		}
	}
}

// ensureDBRecord inserts or updates the Go-native plugin's metadata in the
// plugins table so the admin UI can discover it (settings, pages, frontend assets).
func (m *Manager) ensureDBRecord(ctx context.Context, mf sdk.Manifest) {
	m.ensureDBRecordExt(ctx, mf, "builtin")
}

// isEnabled checks if a plugin is enabled in the DB.
// Returns true if the plugin doesn't exist in DB yet (new Go-native plugin).
func (m *Manager) isEnabled(ctx context.Context, id string) bool {
	total, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).Count()
	if err != nil || total == 0 {
		// Plugin not in DB yet — treat as enabled (new Go-native plugin)
		return true
	}
	// Plugin exists in DB — respect the enabled flag
	enabled, _ := g.DB().Ctx(ctx).Model("plugins").Where("id", id).Where("enabled", 1).Count()
	return enabled > 0
}

// runMigrations executes pending migrations for a plugin.
func (m *Manager) runMigrations(ctx context.Context, pluginID string, migrations []sdk.Migration) error {
	db := g.DB()

	// Get current migration version
	val, _ := db.Ctx(ctx).Model("plugin_migrations").
		Where("plugin_id", pluginID).
		OrderDesc("version").
		Limit(1).
		Value("version")
	currentVersion := val.Int()

	for _, mig := range migrations {
		if mig.Version <= currentVersion {
			continue
		}
		// Resolve dialect-specific SQL
		resolvedUp, err := resolveSQL(mig.Up)
		if err != nil {
			return fmt.Errorf("migration v%d: %w", mig.Version, err)
		}
		// Execute migration statements (may be semicolon-separated)
		for _, stmt := range splitSQL(resolvedUp) {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("migration v%d: %w", mig.Version, err)
			}
		}
		// Record migration
		_, err = db.Exec(ctx,
			"INSERT INTO plugin_migrations (plugin_id, version) VALUES (?, ?)",
			pluginID, mig.Version)
		if err != nil {
			return fmt.Errorf("recording migration v%d: %w", mig.Version, err)
		}
		g.Log().Infof(ctx, "[pluginmgr] %s: applied migration v%d", pluginID, mig.Version)
	}
	return nil
}

// splitSQL splits a string by semicolons, respecting that some statements
// may not use semicolons at all.
func splitSQL(sql string) []string {
	parts := strings.Split(sql, ";")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}
