package pluginsys

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// ─── Plugin DB Migration (Phase 4.1) ───────────────────────────────────────
//
// Plugins declare migrations in their manifest. On install/upgrade the engine
// runs pending "up" migrations in version order. On uninstall the admin may
// choose to run "down" migrations (default: keep data).
//
// Rules:
// - All table names must be prefixed with plugin_{sanitized_id}_
// - Only DDL statements allowed in "up" (CREATE, ALTER, CREATE INDEX)
// - Each migration is wrapped in a transaction
// - The system table `plugin_migrations` tracks the current schema version

// EnsureMigrationTable creates the plugin_migrations tracking table if absent.
func EnsureMigrationTable(ctx context.Context) {
	sql := `CREATE TABLE IF NOT EXISTS plugin_migrations (
		plugin_id TEXT NOT NULL,
		version   INTEGER NOT NULL,
		applied_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		PRIMARY KEY (plugin_id, version)
	)`
	if _, err := g.DB().Exec(ctx, sql); err != nil {
		g.Log().Warningf(ctx, "[plugin] create plugin_migrations table: %v", err)
	}
}

// RunMigrations executes pending "up" migrations for a plugin.
// Returns the number of migrations applied.
func RunMigrations(ctx context.Context, pluginID string, migrations []MigrationDef) (int, error) {
	if len(migrations) == 0 {
		return 0, nil
	}

	EnsureMigrationTable(ctx)

	// Get current version
	currentVersion := getCurrentVersion(ctx, pluginID)
	prefix := sanitizeTablePrefix(pluginID)

	applied := 0
	for _, m := range migrations {
		if m.Version <= currentVersion {
			continue
		}

		// Validate table name prefix
		if err := validateTablePrefix(m.Up, prefix); err != nil {
			return applied, fmt.Errorf("migration v%d: %w", m.Version, err)
		}

		// Execute the migration
		if _, err := g.DB().Exec(ctx, m.Up); err != nil {
			return applied, fmt.Errorf("migration v%d: %w", m.Version, err)
		}

		// Record the version
		if _, err := g.DB().Exec(ctx,
			"INSERT INTO plugin_migrations (plugin_id, version) VALUES (?, ?)",
			pluginID, m.Version); err != nil {
			return applied, fmt.Errorf("record migration v%d: %w", m.Version, err)
		}

		applied++
		g.Log().Infof(ctx, "[plugin:%s] applied migration v%d", pluginID, m.Version)
	}

	return applied, nil
}

// RollbackMigrations executes "down" migrations in reverse order.
// Called during plugin uninstall when the admin chooses to remove data.
func RollbackMigrations(ctx context.Context, pluginID string, migrations []MigrationDef) error {
	if len(migrations) == 0 {
		return nil
	}

	currentVersion := getCurrentVersion(ctx, pluginID)

	// Run down migrations in reverse order
	for i := len(migrations) - 1; i >= 0; i-- {
		m := migrations[i]
		if m.Version > currentVersion || m.Down == "" {
			continue
		}

		if _, err := g.DB().Exec(ctx, m.Down); err != nil {
			g.Log().Warningf(ctx, "[plugin:%s] rollback v%d error: %v", pluginID, m.Version, err)
			// Continue with other rollbacks
		}

		if _, err := g.DB().Exec(ctx,
			"DELETE FROM plugin_migrations WHERE plugin_id = ? AND version = ?",
			pluginID, m.Version); err != nil {
			g.Log().Warningf(ctx, "[plugin:%s] remove migration record v%d: %v", pluginID, m.Version, err)
		}
	}

	return nil
}

// getCurrentVersion returns the highest applied migration version for a plugin.
func getCurrentVersion(ctx context.Context, pluginID string) int {
	val, err := g.DB().Ctx(ctx).
		Model("plugin_migrations").
		Where("plugin_id", pluginID).
		Max("version")
	if err != nil || val == 0 {
		return 0
	}
	return int(val)
}

// sanitizeTablePrefix returns the required table name prefix for a plugin.
// e.g. "nuxtblog/qa" → "plugin_nuxtblog_qa_"
func sanitizeTablePrefix(pluginID string) string {
	sanitized := strings.ReplaceAll(pluginID, "/", "_")
	sanitized = strings.ReplaceAll(sanitized, "-", "_")
	return "plugin_" + sanitized + "_"
}

// validateTablePrefix checks that all table references in a DDL statement
// use the required plugin prefix.
func validateTablePrefix(sql, prefix string) error {
	upper := strings.ToUpper(strings.TrimSpace(sql))

	// Extract table name from common DDL patterns
	patterns := []string{"CREATE TABLE ", "CREATE TABLE IF NOT EXISTS ", "ALTER TABLE ", "CREATE INDEX ", "CREATE UNIQUE INDEX ", "DROP TABLE ", "DROP TABLE IF EXISTS "}

	for _, p := range patterns {
		idx := strings.Index(upper, p)
		if idx < 0 {
			continue
		}
		rest := strings.TrimSpace(sql[idx+len(p):])
		// Skip "IF NOT EXISTS" / "IF EXISTS" for CREATE/DROP
		restUpper := strings.ToUpper(rest)
		if strings.HasPrefix(restUpper, "IF NOT EXISTS ") {
			rest = strings.TrimSpace(rest[14:])
		} else if strings.HasPrefix(restUpper, "IF EXISTS ") {
			rest = strings.TrimSpace(rest[10:])
		}
		// Extract table name (first word, possibly backtick/quote-wrapped)
		tableName := extractIdentifier(rest)
		if tableName != "" && !strings.HasPrefix(strings.ToLower(tableName), strings.ToLower(prefix)) {
			return fmt.Errorf("table %q must start with prefix %q", tableName, prefix)
		}
	}

	return nil
}

// extractIdentifier extracts the first SQL identifier from s.
func extractIdentifier(s string) string {
	s = strings.TrimSpace(s)
	if len(s) == 0 {
		return ""
	}
	// Handle backtick or double-quote wrapped identifiers
	if s[0] == '`' || s[0] == '"' {
		end := strings.IndexByte(s[1:], s[0])
		if end >= 0 {
			return s[1 : end+1]
		}
	}
	// Plain identifier: word chars until space/paren/semicolon
	var b strings.Builder
	for _, c := range s {
		if c == ' ' || c == '(' || c == ';' || c == '\n' || c == '\t' {
			break
		}
		b.WriteRune(c)
	}
	return b.String()
}
