package cmd

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	medialogic "github.com/nuxtblog/nuxtblog/internal/logic/media"
	"github.com/nuxtblog/nuxtblog/internal/util/password"
	dbsql "github.com/nuxtblog/nuxtblog/sql"
)

// autoMigrate creates the database schema (if needed) and seeds required
// default data on every startup. It is safe to run on every boot — all
// CREATE TABLE statements use IF NOT EXISTS and data inserts are guarded
// by existence checks.
func autoMigrate(ctx context.Context) error {
	db := g.DB()

	// ── 1. Auto-create schema ─────────────────────────────────────────────
	// Detect the configured database type and execute the matching schema.
	// All statements use IF NOT EXISTS so this is idempotent.
	dbType := strings.ToLower(db.GetConfig().Type)
	var schema string
	switch dbType {
	case "pgsql", "postgres", "postgresql":
		schema = dbsql.PostgreSQL
	default: // sqlite and anything else
		schema = dbsql.SQLite
	}
	for _, stmt := range splitSQLStatements(schema) {
		if _, err := db.Ctx(ctx).Exec(ctx, stmt); err != nil {
			// Log at debug level — most "errors" are harmless (object already exists)
			g.Log().Debugf(ctx, "[schema] %v", err)
		}
	}

	// ── 2. Seed default options ───────────────────────────────────────────
	// Values are only inserted if the key does not yet exist.
	defaultOpts := []struct{ key, value string }{
		// ── Roles ─────────────────────────────────────────────────────────
		{"custom_role_defs", `[]`},
		// Empty object = use frontend code defaults (permissions.ts).
		// Only populate with admin-defined overrides.
		{"role_capabilities", `{}`},
		// ── General site settings ──────────────────────────────────────────
		{"site_name", `"My Blog"`},
		{"site_description", `""`},
		{"site_url", `"http://localhost:3000"`},
		{"admin_email", `""`},
		{"language", `"zh-CN"`},
		{"allow_registration", `false`},
		{"default_post_cover", `"/images/default-cover.svg"`},
		{"error_post_cover", `"/images/default-cover.svg"`},
		{"default_avatar", `"/images/default-avatar.svg"`},
		{"default_user_bg", `"/images/default-user-banner.jpg"`},
		{"default_avatar_type", `"image"`},
		{"default_avatar_url", `"/images/default-avatar.svg"`},
		{"footer_text", `""`},
		{"icp_number", `""`},
		{"police_number", `""`},
		// ── Reading ───────────────────────────────────────────────────────
		{"posts_per_page", `10`},
		// ── Discussion ────────────────────────────────────────────────────
		{"default_allow_comments", `true`},
		{"comment_moderation", `false`},
		{"comment_require_name_email", `true`},
		{"comment_max_links", `2`},
		{"comment_blacklist", `""`},
		// ── Writing ───────────────────────────────────────────────────────
		{"default_editor", `"markdown"`},
		{"auto_save", `true`},
		{"auto_save_interval", `60`},
		// ── Homepage ──────────────────────────────────────────────────────
		// All widgets and sections enabled by default so new installs show the full feature set.
		{"homepage_sidebar_enabled", `true`},
		{"homepage_sidebar_widgets", `[` +
			`{"id":"user_box","label":"admin.widgets.id_user_box","enabled":true},` +
			`{"id":"search","label":"admin.widgets.id_search","enabled":true,"showRecent":true,"showHot":true},` +
			`{"id":"tags","label":"admin.widgets.id_tags","enabled":true,"maxCount":15},` +
			`{"id":"latest_posts","label":"admin.widgets.id_latest_posts","enabled":true,"maxCount":5},` +
			`{"id":"latest_comments","label":"admin.widgets.id_latest_comments","enabled":true,"maxCount":5},` +
			`{"id":"recommend","label":"admin.widgets.id_recommend","enabled":true,"maxCount":5},` +
			`{"id":"featured","label":"admin.widgets.id_featured","enabled":true,"maxCount":5},` +
			`{"id":"random_posts","label":"admin.widgets.id_random_posts","enabled":true,"maxCount":5}` +
			`]`},
		{"homepage_sections", `[` +
			`{"id":"latest","label":"admin.settings.homepage.section_latest","enabled":true,"count":6,"layout":"grid","includeCategoryIds":[],"excludeCategoryIds":[]},` +
			`{"id":"hot","label":"admin.settings.homepage.section_hot","enabled":true,"count":8,"layout":"ranking"},` +
			`{"id":"featured","label":"admin.settings.homepage.section_featured","enabled":true,"count":6,"layout":"hero"},` +
			`{"id":"random","label":"admin.settings.homepage.section_random","enabled":true,"count":8,"layout":"grid"},` +
			`{"id":"timeline","label":"admin.settings.homepage.section_timeline","enabled":true,"count":10,"layout":"timeline"},` +
			`{"id":"masonry","label":"admin.settings.homepage.section_masonry","enabled":true,"count":9,"layout":"masonry"}` +
			`]`},
		// ── Article (blog/post) sidebar ────────────────────────────────────
		{"post_sidebar_enabled", `true`},
		{"blog_sidebar_widgets", `[` +
			`{"id":"toc","label":"admin.widgets.id_toc","enabled":true},` +
			`{"id":"author","label":"admin.widgets.id_author","enabled":true},` +
			`{"id":"search","label":"admin.widgets.id_search","enabled":true,"showRecent":true,"showHot":true},` +
			`{"id":"tags","label":"admin.widgets.id_tags","enabled":true,"maxCount":15},` +
			`{"id":"latest_posts","label":"admin.widgets.id_latest_posts","enabled":true,"maxCount":5},` +
			`{"id":"latest_comments","label":"admin.widgets.id_latest_comments","enabled":true,"maxCount":5},` +
			`{"id":"recommend","label":"admin.widgets.id_recommend","enabled":true,"maxCount":5},` +
			`{"id":"featured","label":"admin.widgets.id_featured","enabled":true,"maxCount":5},` +
			`{"id":"random_posts","label":"admin.widgets.id_random_posts","enabled":true,"maxCount":5},` +
			`{"id":"downloads","label":"admin.widgets.id_downloads","enabled":true}` +
			`]`},
		// ── Page sidebar ───────────────────────────────────────────────────
		{"page_sidebar_enabled", `true`},
		{"page_sidebar_widgets", `[` +
			`{"id":"author","label":"admin.widgets.id_author","enabled":true},` +
			`{"id":"search","label":"admin.widgets.id_search","enabled":true,"showRecent":true,"showHot":true},` +
			`{"id":"tags","label":"admin.widgets.id_tags","enabled":true,"maxCount":15},` +
			`{"id":"latest_posts","label":"admin.widgets.id_latest_posts","enabled":true,"maxCount":5},` +
			`{"id":"latest_comments","label":"admin.widgets.id_latest_comments","enabled":true,"maxCount":5},` +
			`{"id":"recommend","label":"admin.widgets.id_recommend","enabled":true,"maxCount":5},` +
			`{"id":"featured","label":"admin.widgets.id_featured","enabled":true,"maxCount":5},` +
			`{"id":"random_posts","label":"admin.widgets.id_random_posts","enabled":true,"maxCount":5}` +
			`]`},
		// ── Media ─────────────────────────────────────────────────────────
		{"media_upload_path", `"{year}/{month}"`},
		// ── Notifications & auth ──────────────────────────────────────────
		{"notify_email", `{"host":"","port":587,"username":"","password":"","from":"Blog <noreply@example.com>","site_name":"My Blog","site_url":"http://localhost:3000"}`},
		{"notify_sms", `{"provider":"","access_key_id":"","access_key_secret":"","sign_name":"","template_code":""}`},
		{"notify_webhook", `{"urls":[],"secret":""}`},
		{"auth_register_verify", `{"mode":"none"}`},
		// ── OAuth providers (DB-managed, overrides config.yaml) ────────────
		{"oauth_github", `{"enabled":false,"clientId":"","clientSecret":"","callbackUrl":"http://localhost:9000/api/v1/auth/oauth/github/callback"}`},
		{"oauth_google", `{"enabled":false,"clientId":"","clientSecret":"","callbackUrl":"http://localhost:9000/api/v1/auth/oauth/google/callback"}`},
		{"oauth_qq", `{"enabled":false,"clientId":"","clientSecret":"","callbackUrl":"http://localhost:9000/api/v1/auth/oauth/qq/callback"}`},
		// oauth_providers: list of custom (frontend-added) provider slugs
		{"oauth_providers", `[]`},
	}
	for _, opt := range defaultOpts {
		cnt, _ := db.Ctx(ctx).Model("options").Where("key", opt.key).Count()
		if cnt == 0 {
			_, _ = db.Ctx(ctx).Model("options").Data(g.Map{
				"key":      opt.key,
				"value":    opt.value,
				"autoload": 1,
			}).Insert()
		}
	}

	// ── 3. Migrate stale role_capabilities ───────────────────────────────
	// If the DB has the old capability names (pre-permissions.ts refactor) the
	// frontend's can('access_admin') check will always fail, locking everyone
	// out of the admin panel.  Detect this by checking whether role 3's cap
	// list contains "access_admin"; if not, reset to {} so the frontend falls
	// back to its hardcoded defaults in permissions.ts.
	if err := fixStaleRoleCapabilities(ctx); err != nil {
		g.Log().Warningf(ctx, "fixStaleRoleCapabilities: %v", err)
	}

	// ── 3b. Add settings_schema column to plugins (idempotent) ───────────
	_, _ = db.Ctx(ctx).Exec(ctx, `ALTER TABLE plugins ADD COLUMN settings_schema TEXT NOT NULL DEFAULT '[]'`)

	// ── 3b3. Add styles column to plugins (idempotent) ───────────────────
	_, _ = db.Ctx(ctx).Exec(ctx, `ALTER TABLE plugins ADD COLUMN styles TEXT NOT NULL DEFAULT ''`)

	// ── 3b4. Add priority column to plugins (idempotent) ─────────────────
	_, _ = db.Ctx(ctx).Exec(ctx, `ALTER TABLE plugins ADD COLUMN priority INTEGER NOT NULL DEFAULT 10`)

	// ── 3b5. Add capabilities column to plugins (idempotent) ─────────────
	_, _ = db.Ctx(ctx).Exec(ctx, `ALTER TABLE plugins ADD COLUMN capabilities TEXT NOT NULL DEFAULT '{}'`)

	// ── 3b2. Add email_verified column to users (idempotent) ─────────────
	_, _ = db.Ctx(ctx).Exec(ctx, `ALTER TABLE users ADD COLUMN email_verified INTEGER NOT NULL DEFAULT 0`)

	// ── 3c. Migrate user_likes / user_bookmarks to polymorphic schema ─────
	// Old schema: PRIMARY KEY (user_id, post_id) — hardcoded to posts table.
	// New schema: PRIMARY KEY (user_id, object_type, object_id) — covers post | doc | moment.
	// Existing rows are preserved: post_id → object_id, object_type = 'post'.
	migrateUserInteractionTables(ctx)

	// ── 4. Sync built-in media categories ────────────────────────────────
	if err := medialogic.SyncBuiltinCategories(ctx); err != nil {
		g.Log().Warningf(ctx, "SyncBuiltinCategories warning: %v", err)
	}

	// ── 5. Seed default admin user ────────────────────────────────────────
	count, err := db.Ctx(ctx).Model("users").Where("id", 1).WhereNull("deleted_at").Count()
	if err != nil {
		return err
	}
	if count == 0 {
		const defaultPassword = "admin123"
		hash, err := password.Hash(defaultPassword)
		if err != nil {
			return err
		}
		now := gtime.Now()
		_, err = db.Ctx(ctx).Model("users").Data(g.Map{
			"username":       "admin",
			"email":          "admin@example.com",
			"password_hash":  string(hash),
			"display_name":   "Administrator",
			"role":           3,
			"status":         1,
			"email_verified": 1,
			"locale":         "zh-CN",
			"bio":            "",
			"created_at":     now,
			"updated_at":     now,
		}).Insert()
		if err != nil {
			return err
		}
		g.Log().Infof(ctx, "✅ Default admin created — username: admin  password: %s  (change this immediately)", defaultPassword)
	}

	return nil
}

// migrateUserInteractionTables converts user_likes and user_bookmarks from the
// old hardcoded (user_id, post_id) schema to the polymorphic
// (user_id, object_type, object_id) schema.  Existing rows survive: post_id
// is mapped to object_id with object_type='post'.  Idempotent: skips when
// the object_type column already exists.
func migrateUserInteractionTables(ctx context.Context) {
	db := g.DB()
	dbType := strings.ToLower(db.GetConfig().Type)
	isPG := dbType == "pgsql" || dbType == "postgres" || dbType == "postgresql"

	for _, tbl := range []string{"user_likes", "user_bookmarks"} {
		// Check whether already migrated by looking for the object_type column.
		var isMigrated bool
		if isPG {
			res, err := db.GetAll(ctx,
				"SELECT COUNT(*) AS c FROM information_schema.columns WHERE table_name=$1 AND column_name='object_type'", tbl)
			isMigrated = err == nil && len(res) > 0 && res[0]["c"].Int() > 0
		} else {
			res, err := db.GetAll(ctx, "PRAGMA table_info("+tbl+")")
			if err == nil {
				for _, row := range res {
					if row["name"].String() == "object_type" {
						isMigrated = true
						break
					}
				}
			}
		}
		if isMigrated {
			continue
		}

		// Recreate the table: rename old → create new → copy data → drop old.
		// Neither table is referenced by FK from any other table, so this is safe.
		tmp := tbl + "_pre_v2"
		var idType, tsType, tsDefault string
		if isPG {
			idType, tsType, tsDefault = "BIGINT", "TIMESTAMPTZ", "NOW()"
		} else {
			idType, tsType, tsDefault = "INTEGER", "DATETIME", "(datetime('now'))"
		}
		stmts := []string{
			"ALTER TABLE " + tbl + " RENAME TO " + tmp,
			"CREATE TABLE " + tbl + " (" +
				"user_id " + idType + " NOT NULL, " +
				"object_type TEXT NOT NULL DEFAULT 'post', " +
				"object_id " + idType + " NOT NULL, " +
				"created_at " + tsType + " NOT NULL DEFAULT " + tsDefault + ", " +
				"PRIMARY KEY (user_id, object_type, object_id))",
			"INSERT INTO " + tbl + " (user_id, object_type, object_id, created_at) " +
				"SELECT user_id, 'post', post_id, created_at FROM " + tmp,
			"DROP TABLE " + tmp,
			"CREATE INDEX IF NOT EXISTS idx_" + tbl + "_object ON " + tbl + " (object_type, object_id)",
		}
		ok := true
		for _, stmt := range stmts {
			if _, err := db.Ctx(ctx).Exec(ctx, stmt); err != nil {
				g.Log().Warningf(ctx, "[migrate] %s polymorphic migration failed: %v", tbl, err)
				ok = false
				break
			}
		}
		if ok {
			g.Log().Infof(ctx, "✅ Migrated %s → polymorphic (object_type, object_id)", tbl)
		}
	}
}

// fixStaleRoleCapabilities resets role_capabilities to {} when it was seeded
// with the old pre-refactor capability names that no longer exist in the
// frontend's permissions.ts (e.g. "read_posts", "manage_profile").
// After reset the frontend falls back to its hardcoded DEFAULT_ROLE_CAPABILITIES
// which always contain "access_admin" for roles 2, 3, and 4.
func fixStaleRoleCapabilities(ctx context.Context) error {
	val, err := g.DB().Ctx(ctx).Model("options").Where("key", "role_capabilities").Value("value")
	if err != nil || val.IsNil() {
		return err
	}

	// Value may be JSON-encoded twice (string wrapping an object string)
	raw := val.String()
	var payload string
	if json.Unmarshal([]byte(raw), &payload) == nil {
		raw = payload
	}

	var caps map[string][]string
	if err := json.Unmarshal([]byte(raw), &caps); err != nil {
		return nil // unparseable — leave as-is
	}

	// Consider data stale if role 3 has entries but none is "access_admin"
	role3 := caps["3"]
	if len(role3) == 0 {
		return nil // already empty or no override, nothing to fix
	}
	for _, c := range role3 {
		if c == "access_admin" {
			return nil // already up-to-date
		}
	}

	_, err = g.DB().Ctx(ctx).Model("options").
		Where("key", "role_capabilities").
		Data(g.Map{"value": "{}"}).
		Update()
	if err == nil {
		g.Log().Infof(ctx, "✅ Migrated stale role_capabilities to {} (frontend defaults now apply)")
	}
	return err
}

// splitSQLStatements splits a SQL script into individual executable statements.
// It handles:
//   - Dollar-quoted blocks ($$...$$) used in PostgreSQL PL/pgSQL functions
//   - Single-quoted string literals ('...')
//   - Line comments (--)
//   - Statement terminator (;)
func splitSQLStatements(script string) []string {
	var stmts []string
	var buf strings.Builder
	i, n := 0, len(script)

	for i < n {
		// Line comment: skip to end of line
		if i+1 < n && script[i] == '-' && script[i+1] == '-' {
			for i < n && script[i] != '\n' {
				buf.WriteByte(script[i])
				i++
			}
			continue
		}

		// Dollar-quoted block: $$...$$  (PostgreSQL function bodies)
		if i+1 < n && script[i] == '$' && script[i+1] == '$' {
			buf.WriteByte(script[i])
			buf.WriteByte(script[i+1])
			i += 2
			for i+1 < n {
				if script[i] == '$' && script[i+1] == '$' {
					buf.WriteByte(script[i])
					buf.WriteByte(script[i+1])
					i += 2
					break
				}
				buf.WriteByte(script[i])
				i++
			}
			continue
		}

		// Single-quoted string literal: '...' (handle '' escape)
		if script[i] == '\'' {
			buf.WriteByte(script[i])
			i++
			for i < n {
				if script[i] == '\'' {
					buf.WriteByte(script[i])
					i++
					if i < n && script[i] == '\'' {
						buf.WriteByte(script[i])
						i++
					} else {
						break
					}
				} else {
					buf.WriteByte(script[i])
					i++
				}
			}
			continue
		}

		// Statement terminator
		if script[i] == ';' {
			if stmt := strings.TrimSpace(buf.String()); stmt != "" {
				stmts = append(stmts, stmt)
			}
			buf.Reset()
			i++
			continue
		}

		buf.WriteByte(script[i])
		i++
	}
	// Flush any trailing content (no terminating semicolon)
	if stmt := strings.TrimSpace(buf.String()); stmt != "" {
		stmts = append(stmts, stmt)
	}
	return stmts
}
