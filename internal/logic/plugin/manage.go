package plugin

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/plugin/v1"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ── GetSettings ───────────────────────────────────────────────────────────

func (s *sPlugin) GetSettings(ctx context.Context, id string) (*v1.PluginGetSettingsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type row struct {
		SettingsSchema string `orm:"settings_schema"`
		Settings       string `orm:"settings"`
	}
	var r row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Where("id", id).
		Fields("settings_schema, settings").
		Scan(&r); err != nil {
		return nil, err
	}

	var schema []v1.SettingField
	if r.SettingsSchema != "" && r.SettingsSchema != "[]" {
		_ = json.Unmarshal([]byte(r.SettingsSchema), &schema)
	}
	if schema == nil {
		schema = []v1.SettingField{}
	}

	var values map[string]interface{}
	if r.Settings != "" && r.Settings != "{}" {
		_ = json.Unmarshal([]byte(r.Settings), &values)
	}
	if values == nil {
		values = map[string]interface{}{}
	}

	return &v1.PluginGetSettingsRes{Schema: schema, Values: values}, nil
}

// ── UpdateSettings ────────────────────────────────────────────────────────

func (s *sPlugin) UpdateSettings(ctx context.Context, id string, values map[string]interface{}) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	raw, err := json.Marshal(values)
	if err != nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "invalid values")
	}
	_, err = g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Data(g.Map{"settings": string(raw), "updated_at": gtime.Now()}).Update()
	return err
}

// ── Uninstall ─────────────────────────────────────────────────────────────

func (s *sPlugin) Uninstall(ctx context.Context, id string) (bool, error) {
	if err := requireAdmin(ctx); err != nil {
		return false, err
	}

	// Look up source before deleting the row so we know whether this is a
	// builtin plugin that needs symmetric cleanup of builtin/{dirName}/.
	srcVal, _ := g.DB().Ctx(ctx).Model("plugins").Where("id", id).Value("source")
	src := srcVal.String()
	isBuiltin := src == "builtin"

	// Unload from Goja engine (legacy JS plugins)
	eng.Unload(id)

	// Unload from plugin manager (JS/builtin plugins) — cascade unloads dependents
	if mgr := eng.GetManager(); mgr != nil {
		unloaded := mgr.UnloadPlugin(id)
		if len(unloaded) > 1 {
			g.Log().Infof(ctx, "[plugin] cascade unloaded %d plugins: %v", len(unloaded), unloaded)
		}
		// Fully uninstall cascade-unloaded dependents (not just disable).
		for _, pid := range unloaded {
			if pid == id {
				continue
			}
			// 1. Remove plugin directory
			depDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(pid))
			if info, err := os.Stat(depDir); err == nil && info.IsDir() {
				if rmErr := os.RemoveAll(depDir); rmErr != nil {
					g.Log().Warningf(ctx, "[plugin] failed to remove cascaded plugin dir %s: %v", depDir, rmErr)
				}
			}
			// 2. Rollback migrations
			s.rollbackPluginMigrations(ctx, pid)
			// 3. Remove KV store entries
			_, _ = g.DB().Ctx(ctx).Model("options").
				WhereLike("key", "plugin_store:"+pid+":%").Delete()
			// 4. Delete DB record
			if _, err := g.DB().Ctx(ctx).Model("plugins").Where("id", pid).Delete(); err != nil {
				g.Log().Warningf(ctx, "[plugin] delete cascaded dependent %s failed: %v", pid, err)
			}
		}
	}

	// Remove plugin directory from data/plugins/{id}/
	pluginDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(id))
	if info, err := os.Stat(pluginDir); err == nil && info.IsDir() {
		if rmErr := os.RemoveAll(pluginDir); rmErr != nil {
			g.Log().Warningf(ctx, "[plugin] failed to remove plugin dir %s: %v", pluginDir, rmErr)
		}
	}

	// Rollback database migrations (drop plugin-created tables)
	s.rollbackPluginMigrations(ctx, id)

	// Also remove per-plugin KV store entries
	_, _ = g.DB().Ctx(ctx).Model("options").
		WhereLike("key", "plugin_store:"+id+":%").Delete()

	// Symmetric cleanup for builtin plugins: remove the source tree under
	// builtin/{dirName}/ and regenerate builtin/plugins.go so the next build
	// no longer compiles this plugin in. The currently running binary still
	// contains the compiled code — the warning tells the user to rebuild.
	if src == "builtin" {
		dirName := strings.ReplaceAll(id, "-", "_") // mirrors installBuiltinPlugin
		builtinDir := builtinSourceDir()
		if builtinDir != "" {
			target := filepath.Join(builtinDir, dirName)
			if info, err := os.Stat(target); err == nil && info.IsDir() {
				if rmErr := os.RemoveAll(target); rmErr != nil {
					g.Log().Warningf(ctx, "[plugin] failed to remove builtin source %s: %v", target, rmErr)
				}
			}
			if regenErr := regenerateBuiltinImports(builtinDir); regenErr != nil {
				g.Log().Warningf(ctx, "[plugin] failed to regenerate builtin imports: %v", regenErr)
			}
			g.Log().Warningf(ctx, "[plugin] builtin plugin %s source removed; rebuild and restart required for the compiled code to disappear", id)
		}
	}

	_, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).Delete()
	if err != nil {
		return false, err
	}
	return isBuiltin, nil
}

// rollbackPluginMigrations reads the plugin's manifest from the DB and runs
// down migrations to drop plugin-created tables, then cleans up the
// plugin_migrations tracking records.
func (s *sPlugin) rollbackPluginMigrations(ctx context.Context, id string) {
	manifestVal, _ := g.DB().Ctx(ctx).Model("plugins").
		Where("id", id).Value("manifest")
	manifestStr := manifestVal.String()
	if manifestStr == "" || manifestStr == "{}" {
		return
	}

	// Try to extract migrations from the stored manifest JSON.
	// The manifest may be either sdk.Manifest format (yaml-based) or
	// pluginsys.Manifest format (package.json-based). Both have a
	// "migrations" array with "version", "up", "down" fields.
	var mf struct {
		Migrations []eng.MigrationDef `json:"migrations"`
	}
	if err := json.Unmarshal([]byte(manifestStr), &mf); err != nil || len(mf.Migrations) == 0 {
		// No migrations to roll back — just clean up tracking records
		_, _ = g.DB().Ctx(ctx).Model("plugin_migrations").
			Where("plugin_id", id).Delete()
		return
	}

	if err := eng.RollbackMigrations(ctx, id, mf.Migrations); err != nil {
		g.Log().Warningf(ctx, "[plugin] %s migration rollback error: %v", id, err)
	} else {
		g.Log().Infof(ctx, "[plugin] %s: rolled back %d migration(s)", id, len(mf.Migrations))
	}
}

// ── Toggle ────────────────────────────────────────────────────────────────

func (s *sPlugin) Toggle(ctx context.Context, id string, enabled bool) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}

	enabledInt := 0
	if enabled {
		enabledInt = 1
	}
	_, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Data(g.Map{"enabled": enabledInt, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}

	if enabled {
		// Try to reload as JS (Goja) plugin first (check for plugin dir on disk)
		pluginDir := filepath.Join(PluginAssetsDir(), sanitizePluginPath(id))
		if mgr := eng.GetManager(); mgr != nil {
			if info, statErr := os.Stat(filepath.Join(pluginDir, "plugin.js")); statErr == nil && !info.IsDir() {
				// This is a JS (Goja) plugin — reload from disk
				_ = mgr.InstallJSPlugin(ctx, pluginDir, "plugin.js")
				return nil
			}
		}

		// Fallback: reload legacy script from DB
		type row struct {
			Script       string `orm:"script"`
			Priority     int    `orm:"priority"`
			Capabilities string `orm:"capabilities"`
			Manifest     string `orm:"manifest"`
		}
		var r row
		if e := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
			Fields("script, priority, COALESCE(capabilities,'{}') as capabilities, COALESCE(manifest,'{}') as manifest").Scan(&r); e == nil && r.Script != "" {
			var caps eng.PluginCapabilities
			_ = json.Unmarshal([]byte(r.Capabilities), &caps)
			var mf eng.Manifest
			_ = json.Unmarshal([]byte(r.Manifest), &mf)
			mf.Priority = r.Priority
			mf.Capabilities = caps
			_ = eng.Load(id, r.Script, mf)
		}
	} else {
		// Unload from both engines — cascade unloads dependents
		normalized, _ := eng.NormalizePluginID(id)
		eng.Unload(normalized)
		if mgr := eng.GetManager(); mgr != nil {
			unloaded := mgr.UnloadPlugin(id)
			if len(unloaded) > 1 {
				g.Log().Infof(ctx, "[plugin] cascade unloaded %d plugins on disable: %v", len(unloaded), unloaded)
			}
			for _, pid := range unloaded {
				if pid == id {
					continue
				}
				if _, err := g.DB().Ctx(ctx).Model("plugins").Where("id", pid).
					Data(g.Map{"enabled": 0, "updated_at": gtime.Now()}).Update(); err != nil {
					g.Log().Warningf(ctx, "[plugin] disable cascaded dependent %s failed: %v", pid, err)
				}
			}
		}
	}
	return nil
}

// ── Unload Impact ─────────────────────────────────────────────────────────

func (s *sPlugin) UnloadImpact(ctx context.Context, id string) (*v1.PluginUnloadImpactRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	var willUnload []string
	if mgr := eng.GetManager(); mgr != nil {
		willUnload = mgr.UnloadImpact(id)
	}
	if willUnload == nil {
		willUnload = []string{}
	}
	return &v1.PluginUnloadImpactRes{WillUnload: willUnload}, nil
}

// ── GetStyles ─────────────────────────────────────────────────────────────

// GetStyles returns the concatenated CSS of all enabled plugins.
// This endpoint is public — no auth required.
func (s *sPlugin) GetStyles(ctx context.Context) (string, error) {
	type row struct {
		Styles string `orm:"styles"`
	}
	var rows []row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Fields("styles").Where("enabled", 1).
		Scan(&rows); err != nil {
		return "", err
	}
	var sb strings.Builder
	for _, r := range rows {
		if r.Styles != "" {
			sb.WriteString(r.Styles)
			sb.WriteByte('\n')
		}
	}
	return sb.String(), nil
}

// ── Stats / Errors (4.5-B1 / 4.5-B2 / 4.5-B3) ────────────────────────────

func (s *sPlugin) GetStats(ctx context.Context, id string) (*v1.PluginGetStatsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	snap := eng.GetStats(id)
	if snap == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "plugin not loaded: "+id)
	}
	history := eng.GetHistory(id)
	buckets := make([]v1.WindowBucket, 0, len(history))
	for _, b := range history {
		buckets = append(buckets, v1.WindowBucket{
			Minute:      b.Minute.UTC().Format(time.RFC3339),
			Invocations: b.Invocations,
			Errors:      b.Errors,
		})
	}
	res := &v1.PluginGetStatsRes{
		PluginID:      snap.PluginID,
		Invocations:   snap.Invocations,
		Errors:        snap.Errors,
		AvgDurationMs: snap.AvgDurationMs,
		LastError:     snap.LastError,
		History:       buckets,
	}
	if !snap.LastErrorAt.IsZero() {
		res.LastErrorAt = snap.LastErrorAt.UTC().Format(time.RFC3339)
	}
	return res, nil
}

func (s *sPlugin) GetErrors(ctx context.Context, id string) (*v1.PluginGetErrorsRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	entries := eng.GetErrors(id)
	if entries == nil {
		return nil, gerror.NewCode(gcode.CodeNotFound, "plugin not loaded: "+id)
	}
	items := make([]v1.PluginErrorEntry, 0, len(entries))
	for _, e := range entries {
		items = append(items, v1.PluginErrorEntry{
			At:        e.At.UTC().Format(time.RFC3339),
			EventName: e.EventName,
			Message:   e.Message,
			InputDiff: e.InputDiff,
		})
	}
	return &v1.PluginGetErrorsRes{Items: items}, nil
}

// ── GetManifest / UpdateManifest (P-B11) ─────────────────────────────────────

func (s *sPlugin) GetManifest(ctx context.Context, id string) (*v1.PluginGetManifestRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}
	type row struct {
		Manifest string `orm:"manifest"`
	}
	var r row
	if err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Fields("COALESCE(manifest,'{}') as manifest").Scan(&r); err != nil {
		return nil, err
	}
	return &v1.PluginGetManifestRes{Manifest: r.Manifest}, nil
}

func (s *sPlugin) UpdateManifest(ctx context.Context, id string, manifest string) error {
	if err := requireAdmin(ctx); err != nil {
		return err
	}
	// Validate: must be valid JSON that deserialises to a Manifest
	var mf eng.Manifest
	if err := json.Unmarshal([]byte(manifest), &mf); err != nil {
		return gerror.NewCode(gcode.CodeInvalidParameter, "manifest must be valid JSON: "+err.Error())
	}
	_, err := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Data(g.Map{"manifest": manifest, "updated_at": gtime.Now()}).Update()
	if err != nil {
		return err
	}
	// Reload the plugin in the engine to apply pipeline/webhook changes
	type row struct {
		Script       string `orm:"script"`
		Priority     int    `orm:"priority"`
		Capabilities string `orm:"capabilities"`
	}
	var r row
	if e := g.DB().Ctx(ctx).Model("plugins").Where("id", id).
		Fields("script, priority, COALESCE(capabilities,'{}') as capabilities").Scan(&r); e == nil && r.Script != "" {
		var caps eng.PluginCapabilities
		_ = json.Unmarshal([]byte(r.Capabilities), &caps)
		mf.Priority = r.Priority
		mf.Capabilities = caps
		_ = eng.Load(id, r.Script, mf)
	}
	return nil
}

// ── ClientList (Phase 2.4) ────────────────────────────────────────────────

func (s *sPlugin) ClientList(ctx context.Context) (*v1.PluginClientListRes, error) {
	if err := requireAdmin(ctx); err != nil {
		return nil, err
	}

	type row struct {
		Id       string `orm:"id"`
		Title    string `orm:"title"`
		Icon     string `orm:"icon"`
		Manifest string `orm:"manifest"`
	}
	var rows []row
	if err := g.DB().Ctx(ctx).Model("plugins").
		Where("enabled", 1).
		Fields("id, title, icon, COALESCE(manifest,'{}') as manifest").
		Scan(&rows); err != nil {
		return nil, err
	}

	items := make([]v1.PluginClientItem, 0, len(rows))
	for _, r := range rows {
		var mf eng.Manifest
		_ = json.Unmarshal([]byte(r.Manifest), &mf)

		item := v1.PluginClientItem{
			ID:         r.Id,
			Title:      r.Title,
			Icon:       r.Icon,
			TrustLevel: string(mf.TrustLevel),
			AdminJS:    mf.AdminJS,
			PublicJS:   mf.PublicJS,
		}
		if item.TrustLevel == "" {
			item.TrustLevel = "community"
		}
		if mf.Contributes != nil {
			cb, _ := json.Marshal(mf.Contributes)
			item.Contributes = string(cb)
		}
		if len(mf.Pages) > 0 {
			pb, _ := json.Marshal(mf.Pages)
			item.Pages = string(pb)
		}
		items = append(items, item)
	}
	return &v1.PluginClientListRes{Items: items}, nil
}
