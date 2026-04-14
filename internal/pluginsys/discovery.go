package pluginsys

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gogf/gf/v2/frame/g"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// externalCandidate holds discovery-phase data for an external plugin.
type externalCandidate struct {
	manifest  *sdk.Manifest
	pluginDir string
}

// LoadExternal scans dataDir for plugin directories, reads plugin.yaml,
// and loads JS/full type plugins via Goja.
// Uses two-phase loading: discovery pass, then topological sort, then activate.
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

	// Phase 1: Discovery — parse all manifests
	candidates := make(map[string]*externalCandidate)

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

		// Ensure new plugins on disk are registered in DB (disabled by default).
		// This lets admins see and enable them from the admin panel.
		m.ensureDBRecordDiscovered(ctx, *pm)

		// Check if enabled
		if !m.isEnabled(ctx, pm.ID) {
			g.Log().Infof(ctx, "[pluginmgr] external plugin %s not enabled, skipping", pm.ID)
			continue
		}

		candidates[pm.ID] = &externalCandidate{manifest: pm, pluginDir: pluginDir}

		// Register dependencies for JS/full plugins
		if pm.Type == sdk.TypeJS || pm.Type == sdk.TypeFull {
			m.graph.Add(pm.ID, pm.Depends)
		}
	}

	// Build union of already-loaded (builtin) IDs + all candidate IDs for sorting
	m.mu.RLock()
	candidateIDs := make([]string, 0, len(candidates))
	for id := range candidates {
		candidateIDs = append(candidateIDs, id)
	}
	allIDs := make([]string, 0, len(m.plugins)+len(candidateIDs))
	for id := range m.plugins {
		allIDs = append(allIDs, id)
	}
	m.mu.RUnlock()
	allIDs = append(allIDs, candidateIDs...)

	versionResolver := func(id string) string {
		// Check already-loaded plugins
		m.mu.RLock()
		if lp, ok := m.plugins[id]; ok {
			m.mu.RUnlock()
			return lp.plugin.Manifest().Version
		}
		m.mu.RUnlock()
		// Check candidates
		if c, ok := candidates[id]; ok {
			return c.manifest.Version
		}
		return ""
	}

	sorted, err := m.graph.TopologicalSort(allIDs, versionResolver)
	if err != nil {
		g.Log().Warningf(ctx, "[pluginmgr] dependency sort failed: %v; loading in discovery order", err)
		sorted = allIDs
	}

	// Phase 3: Activate in sorted order, skipping already-loaded builtins
	for _, id := range sorted {
		c, ok := candidates[id]
		if !ok {
			continue // already-loaded builtin, skip
		}
		pm := c.manifest

		switch pm.Type {
		case sdk.TypeJS, sdk.TypeFull:
			jsFile := pm.JSEntry
			if jsFile == "" {
				jsFile = "plugin.js"
			}
			jsPath := filepath.Join(c.pluginDir, jsFile)
			if _, err := os.Stat(jsPath); err != nil {
				g.Log().Warningf(ctx, "[pluginmgr] %s: js source %s not found", pm.ID, jsFile)
				continue
			}
			p, err := loadGojaPlugin(ctx, c.pluginDir, jsFile, *pm)
			if err != nil {
				g.Log().Errorf(ctx, "[pluginmgr] %s goja load failed: %v", pm.ID, err)
				continue
			}
			if err := m.activatePlugin(ctx, p, "external", c.pluginDir, jsFile); err != nil {
				g.Log().Errorf(ctx, "[pluginmgr] %s activation failed: %v", pm.ID, err)
			}
			m.ensureDBRecordExt(ctx, *pm, "external")

		case sdk.TypeBuiltin:
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
