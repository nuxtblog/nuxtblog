// Package pluginsys manages Go-native plugin lifecycle.
//
// It loads statically compiled plugins (via sdk.GetStatic()) and
// dynamically interpreted plugins (via Goja JS) then wires them
// into the platform's route system, event bus, and filter chain.
package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/sdk"
)

// Manager holds all loaded Go-native plugins.
type Manager struct {
	mu      sync.RWMutex
	plugins map[string]*loadedPlugin
	server  *ghttp.Server // set by RegisterRoutes, used for dynamic route binding
}

type loadedPlugin struct {
	plugin  sdk.Plugin
	ctx     sdk.PluginContext
	routes  []routeEntry
	filters []sdk.FilterDef

	// Observability
	stats   *PluginStats
	window  *SlidingWindow
	errors  *ErrorRingBuffer
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
func (m *Manager) LoadStatic(ctx context.Context) error {
	for _, p := range sdk.GetStatic() {
		mf := p.Manifest()
		id := mf.ID
		if id == "" {
			g.Log().Warning(ctx, "[pluginmgr] skipping plugin with empty ID")
			continue
		}

		// Check if plugin is enabled in DB
		if !m.isEnabled(ctx, id) {
			g.Log().Infof(ctx, "[pluginmgr] plugin %s not enabled, skipping", id)
			continue
		}

		if err := m.activatePlugin(ctx, p, "builtin"); err != nil {
			g.Log().Errorf(ctx, "[pluginmgr] plugin %s failed: %v", id, err)
			continue
		}
	}
	return nil
}

// RegisterRoutes registers all Go plugin routes on the GoFrame server.
// It also stores the server reference so that dynamically installed plugins
// can have their routes registered immediately without a restart.
func (m *Manager) RegisterRoutes(s *ghttp.Server) {
	// Allow route overwrite so that plugin updates can re-register routes
	// without causing a fatal duplicate-route error.
	s.SetRouteOverWrite(true)
	m.server = s

	ctx := context.Background()
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id, lp := range m.plugins {
		m.bindPluginRoutes(ctx, s, id, lp)
	}
}

// bindPluginRoutes registers a single plugin's routes on the server.
func (m *Manager) bindPluginRoutes(ctx context.Context, s *ghttp.Server, id string, lp *loadedPlugin) {
	for _, re := range lp.routes {
		handler := m.wrapHandler(re, id)
		pattern := strings.ToUpper(re.method) + ":" + re.path
		s.BindHandler(pattern, handler)
		g.Log().Infof(ctx, "[pluginmgr] registered route %s %s → %s", re.method, re.path, id)
	}
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
			go func(id string, lp *loadedPlugin, ep sdk.HasEvents) {
				start := time.Now()
				ep.OnEvent(ctx, event, data)
				dur := time.Since(start)
				lp.recordExec(id, "event:"+event, dur, nil)
			}(id, lp, ep)
		}
	}
}

// Shutdown deactivates all loaded plugins.
func (m *Manager) Shutdown() {
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
		// Execute migration statements (may be semicolon-separated)
		for _, stmt := range splitSQL(mig.Up) {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if _, err := db.Exec(ctx, stmt); err != nil {
				return fmt.Errorf("migration v%d: %w", mig.Version, err)
			}
		}
		// Record migration
		_, err := db.Exec(ctx,
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

// wrapHandler converts a plugin http.HandlerFunc to a GoFrame handler with auth.
func (m *Manager) wrapHandler(re routeEntry, pluginID string) ghttp.HandlerFunc {
	return func(r *ghttp.Request) {
		// Check if plugin is still loaded (handles dynamic uninstall/disable)
		m.mu.RLock()
		lp, alive := m.plugins[pluginID]
		m.mu.RUnlock()
		if !alive {
			r.Response.WriteStatus(http.StatusNotFound)
			return
		}

		// Parse JWT for user info
		if claims, err := middleware.ParseBearerToken(r); err == nil {
			r.SetCtxVar("user_id", claims.UserID)
			r.SetCtxVar("user_role", claims.Role)
		}

		// Auth check
		uid := r.GetCtxVar("user_id").Int()
		role := r.GetCtxVar("user_role").Int()

		switch re.auth {
		case "admin":
			if uid <= 0 || role < 2 {
				r.Response.WriteJsonExit(g.Map{"code": 401, "message": "admin access required"})
				return
			}
		case "user":
			if uid <= 0 {
				r.Response.WriteJsonExit(g.Map{"code": 401, "message": "authentication required"})
				return
			}
		}

		// Record stats for route execution
		start := time.Now()
		re.handler(r.Response.Writer, r.Request)
		dur := time.Since(start)

		routeErr := error(nil)
		if r.Response.Status >= 500 {
			routeErr = fmt.Errorf("HTTP %d on %s %s", r.Response.Status, re.method, re.path)
		}
		lp.recordExec(pluginID, "route:"+re.path, dur, routeErr)
	}
}

// recordExec records a single execution (route/filter/event) into the plugin's
// observability counters.
func (lp *loadedPlugin) recordExec(pluginID, eventName string, dur time.Duration, err error) {
	if lp.stats != nil {
		lp.stats.record(dur, err)
	}
	if lp.window != nil {
		lp.window.record(err != nil)
	}
	if err != nil && lp.errors != nil {
		lp.errors.Add(ErrorEntry{
			At:        time.Now(),
			EventName: eventName,
			Phase:     "route",
			Message:   err.Error(),
			Duration:  dur,
		})
	}
}

// getPluginObs returns the observability data for a plugin by ID.
func (m *Manager) getPluginObs(id string) *loadedPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.plugins[id]
}

// ─── registrar implements sdk.RouteRegistrar ────────────────────────────────

type registrar struct {
	pluginID string
	entries  []routeEntry
}

func (reg *registrar) Handle(method, path string, handler http.HandlerFunc, opts ...sdk.RouteOption) {
	cfg := sdk.ApplyOptions(opts)
	auth := cfg.Auth
	if auth == "" {
		auth = "public"
	}

	// Enforce /api/plugin/ prefix
	if !strings.HasPrefix(path, "/api/plugin/") {
		path = fmt.Sprintf("/api/plugin/%s%s", reg.pluginID, path)
	}

	reg.entries = append(reg.entries, routeEntry{
		method:  strings.ToUpper(method),
		path:    path,
		handler: handler,
		auth:    auth,
	})
}

// ─── Platform service implementations ───────────────────────────────────────

// pluginDB implements sdk.DB
type pluginDB struct {
	pluginID string
	prefix   string
}

func newPluginDB(pluginID string) *pluginDB {
	return &pluginDB{pluginID: pluginID, prefix: sanitizeTablePrefix(pluginID)}
}

func (d *pluginDB) Query(sql string, args ...any) ([]map[string]any, error) {
	ctx := context.Background()
	result, err := g.DB().GetAll(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	rows := make([]map[string]any, 0, len(result))
	for _, r := range result {
		rows = append(rows, r.Map())
	}
	return rows, nil
}

func (d *pluginDB) Execute(sql string, args ...any) (int64, error) {
	ctx := context.Background()
	res, err := g.DB().Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// pluginStore implements sdk.Store
type pluginStore struct {
	pluginID string
}

func newPluginStore(pluginID string) *pluginStore {
	return &pluginStore{pluginID: pluginID}
}

func (s *pluginStore) storeKey(key string) string {
	return "plugin_store:" + s.pluginID + ":" + key
}

func (s *pluginStore) Get(key string) (any, error) {
	val, err := g.DB().Ctx(context.Background()).
		Model("options").Where("key", s.storeKey(key)).Value("value")
	if err != nil {
		return nil, err
	}
	if val.IsNil() {
		return nil, nil
	}
	var v any
	if json.Unmarshal([]byte(val.String()), &v) == nil {
		return v, nil
	}
	return val.String(), nil
}

func (s *pluginStore) Set(key string, value any) error {
	raw, _ := json.Marshal(value)
	return upsertOption(context.Background(), s.storeKey(key), string(raw))
}

func (s *pluginStore) Delete(key string) error {
	_, err := g.DB().Ctx(context.Background()).Model("options").
		Where("key", s.storeKey(key)).Delete()
	return err
}

func (s *pluginStore) Increment(key string, delta ...int64) (int64, error) {
	d := int64(1)
	if len(delta) > 0 {
		d = delta[0]
	}
	ctx := context.Background()
	db := g.DB()
	dbKey := s.storeKey(key)

	dbType := ""
	if cfg := db.GetConfig(); cfg != nil {
		dbType = strings.ToLower(cfg.Type)
	}
	var sql string
	switch dbType {
	case "mysql", "mariadb":
		sql = "INSERT INTO `options` (`key`, `value`, `autoload`) VALUES (?, ?, 0)" +
			" ON DUPLICATE KEY UPDATE `value` = CAST(`value` AS SIGNED) + CAST(VALUES(`value`) AS SIGNED)"
	default:
		sql = "INSERT INTO options (key, value, autoload) VALUES (?, ?, 0)" +
			" ON CONFLICT(key) DO UPDATE SET value = CAST(options.value AS INTEGER) + CAST(excluded.value AS INTEGER)"
	}
	raw, _ := json.Marshal(d)
	if _, err := db.Exec(ctx, sql, dbKey, string(raw)); err != nil {
		return 0, err
	}
	val, _ := db.Ctx(ctx).Model("options").Where("key", dbKey).Value("value")
	var result int64
	if val != nil && !val.IsNil() {
		_ = json.Unmarshal([]byte(val.String()), &result)
	}
	return result, nil
}

func (s *pluginStore) DeletePrefix(prefix string) (int64, error) {
	dbPrefix := "plugin_store:" + s.pluginID + ":" + prefix
	res, err := g.DB().Ctx(context.Background()).Model("options").
		WhereLike("key", dbPrefix+"%").Delete()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// upsertOption writes key=value into the options table atomically.
func upsertOption(ctx context.Context, key, value string) error {
	db := g.DB()
	dbType := ""
	if cfg := db.GetConfig(); cfg != nil {
		dbType = strings.ToLower(cfg.Type)
	}
	var sql string
	switch dbType {
	case "mysql", "mariadb":
		sql = "INSERT INTO `options` (`key`, `value`, `autoload`) VALUES (?, ?, 0)" +
			" ON DUPLICATE KEY UPDATE `value` = VALUES(`value`)"
	default:
		sql = "INSERT INTO options (key, value, autoload) VALUES (?, ?, 0)" +
			" ON CONFLICT(key) DO UPDATE SET value = excluded.value"
	}
	_, err := db.Exec(ctx, sql, key, value)
	return err
}

// pluginSettings implements sdk.Settings
type pluginSettings struct {
	pluginID string
}

func newPluginSettings(pluginID string) *pluginSettings {
	return &pluginSettings{pluginID: pluginID}
}

func (s *pluginSettings) Get(key string) any {
	return s.GetAll()[key]
}

func (s *pluginSettings) GetAll() map[string]any {
	val, _ := g.DB().Ctx(context.Background()).
		Model("plugins").Where("id", s.pluginID).Value("settings")
	var m map[string]any
	if val != nil && !val.IsNil() {
		_ = json.Unmarshal([]byte(val.String()), &m)
	}
	if m == nil {
		m = map[string]any{}
	}
	return m
}

// pluginLogger implements sdk.Logger
type pluginLogger struct {
	pluginID string
}

func newPluginLogger(pluginID string) *pluginLogger {
	return &pluginLogger{pluginID: pluginID}
}

func (l *pluginLogger) Info(msg string) {
	g.Log().Infof(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}

func (l *pluginLogger) Warn(msg string) {
	g.Log().Warningf(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}

func (l *pluginLogger) Error(msg string) {
	g.Log().Errorf(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}

func (l *pluginLogger) Debug(msg string) {
	g.Log().Infof(context.Background(), "[plugin:%s] %s", l.pluginID, msg)
}
