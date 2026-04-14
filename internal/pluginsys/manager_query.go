package pluginsys

import "fmt"

// ─── pluginQuery implements sdk.PluginQuery ──────────────────────────────────

type pluginQuery struct {
	mgr *Manager
}

func (q *pluginQuery) IsAvailable(id string) bool {
	q.mgr.mu.RLock()
	_, ok := q.mgr.plugins[id]
	q.mgr.mu.RUnlock()
	return ok
}

func (q *pluginQuery) GetVersion(id string) string {
	q.mgr.mu.RLock()
	defer q.mgr.mu.RUnlock()
	if lp, ok := q.mgr.plugins[id]; ok {
		return lp.plugin.Manifest().Version
	}
	return ""
}

func (q *pluginQuery) GetSetting(pluginID, key string) (any, error) {
	q.mgr.mu.RLock()
	lp, ok := q.mgr.plugins[pluginID]
	q.mgr.mu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("plugin %s not available", pluginID)
	}
	mf := lp.plugin.Manifest()
	for _, s := range mf.Settings {
		if s.Key == key {
			if !s.Shared {
				return nil, fmt.Errorf("setting %s.%s is not shared", pluginID, key)
			}
			settings := newPluginSettings(pluginID)
			return settings.Get(key), nil
		}
	}
	return nil, fmt.Errorf("setting %s.%s not found", pluginID, key)
}

// UnloadImpact returns plugin IDs that would be cascade-unloaded (excluding id itself).
func (m *Manager) UnloadImpact(id string) []string {
	cascade := m.graph.GetCascadeUnloadOrder(id)
	result := make([]string, 0, len(cascade)-1)
	for _, pid := range cascade {
		if pid != id {
			result = append(result, pid)
		}
	}
	return result
}
