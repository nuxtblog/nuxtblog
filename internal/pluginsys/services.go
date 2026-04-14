package pluginsys

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// ─── Platform service implementations ───────────────────────────────────────

// pluginDB implements sdk.DB
type pluginDB struct {
	pluginID string
	prefix   string
	caps     *DBCap
	trust    TrustLevel
}

func newPluginDB(pluginID string, caps *DBCap, trust TrustLevel) *pluginDB {
	return &pluginDB{
		pluginID: pluginID,
		prefix:   sanitizeTablePrefix(pluginID),
		caps:     caps,
		trust:    trust,
	}
}

func (d *pluginDB) Query(sql string, args ...any) ([]map[string]any, error) {
	if err := (&sqlGuard{d.prefix, d.caps, d.trust}).validate(sql); err != nil {
		return nil, fmt.Errorf("db access denied: %w", err)
	}
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
	if err := (&sqlGuard{d.prefix, d.caps, d.trust}).validate(sql); err != nil {
		return 0, fmt.Errorf("db access denied: %w", err)
	}
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
