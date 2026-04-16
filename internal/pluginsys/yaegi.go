package pluginsys

import (
	"context"
	"encoding/json"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// UnloadPlugin deactivates and removes a plugin and all its dependents (cascade).
// Returns the list of plugin IDs that were actually unloaded.
func (m *Manager) UnloadPlugin(id string) []string {
	cascade := m.graph.GetCascadeUnloadOrder(id)

	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()
	var unloaded []string
	for _, pid := range cascade {
		if lp, ok := m.plugins[pid]; ok {
			if dp, ok := lp.plugin.(sdk.HasDeactivate); ok {
				if err := dp.Deactivate(); err != nil {
					g.Log().Warningf(ctx, "[pluginmgr] %s deactivate error: %v", pid, err)
				}
			}
			unregisterPluginMedia(pid)
			delete(m.plugins, pid)
			m.graph.Remove(pid)
			unloaded = append(unloaded, pid)
		}
	}
	return unloaded
}

// ReactivatePlugin deactivates then re-activates a loaded plugin so that it
// can pick up new settings (e.g., storage credentials added after first boot).
// Returns false if the plugin is not currently loaded.
func (m *Manager) ReactivatePlugin(ctx context.Context, id string) bool {
	m.mu.Lock()
	lp, ok := m.plugins[id]
	if !ok {
		m.mu.Unlock()
		return false
	}
	m.mu.Unlock()

	// 1. Deactivate
	if dp, ok := lp.plugin.(sdk.HasDeactivate); ok {
		_ = dp.Deactivate()
	}
	unregisterPluginMedia(id)

	// 2. Invalidate settings cache so fresh values are loaded
	if ps, ok := lp.ctx.Settings.(*pluginSettings); ok {
		ps.InvalidateCache()
	}

	// 3. Rebuild plugin context with fresh settings from DB
	mf := lp.plugin.Manifest()
	var dbCaps *DBCap
	if mf.Capabilities != nil && mf.Capabilities.DB != nil {
		dbCaps = convertSDKDBCap(mf.Capabilities.DB)
	}
	trust := TrustLevel(mf.TrustLevel)
	source := "external"
	m.mu.RLock()
	if lp.pluginDir == "" && lp.jsFile == "" {
		source = "builtin"
	}
	m.mu.RUnlock()
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
		I18n:     mf.I18n,
	}

	// 4. Re-activate
	if err := lp.plugin.Activate(pctx); err != nil {
		g.Log().Warningf(ctx, "[pluginmgr] reactivate %s failed: %v", id, err)
		return true
	}

	// 5. Re-register media categories from capabilities.media
	if mf.Capabilities != nil && mf.Capabilities.Media != nil {
		for _, cat := range mf.Capabilities.Media.Categories {
			cat.ResolveCategoryLabel(mf.I18n)
			_ = pctx.Media.RegisterCategory(cat)
		}
	}

	m.mu.Lock()
	lp.ctx = pctx
	m.mu.Unlock()

	g.Log().Infof(ctx, "[pluginmgr] reactivated plugin: %s", id)
	return true
}

// convertSDKDBCap converts sdk.DBCapability to pluginsys.DBCap.
func convertSDKDBCap(c *sdk.DBCapability) *DBCap {
	if c == nil {
		return nil
	}
	d := &DBCap{
		Own: c.Own,
		Raw: c.Raw,
	}
	for _, t := range c.Tables {
		d.Tables = append(d.Tables, DBTableRule{
			Table: t.Table,
			Ops:   t.Ops,
		})
	}
	return d
}

// manifestCapsJSON serializes the sdk.Manifest capabilities to JSON for DB storage.
func manifestCapsJSON(mf sdk.Manifest) string {
	caps := map[string]any{
		"store": map[string]any{"read": true, "write": true},
	}
	if mf.Capabilities != nil && mf.Capabilities.DB != nil {
		caps["db"] = mf.Capabilities.DB
	}
	b, err := json.Marshal(caps)
	if err != nil {
		return `{"store":{"read":true,"write":true}}`
	}
	return string(b)
}
