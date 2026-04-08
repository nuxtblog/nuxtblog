// Package permission implements the IPermission service.
// It reads the role_capabilities option from the database (cached for 1 minute)
// and checks whether a given role possesses a named capability.
//
// Design notes:
//   - Custom roles (ID > 4) are resolved to their base system role via middleware.ResolveRole.
//   - When the database has no entry for a role, hardcoded defaults (matching
//     frontend/app/config/permissions.ts) are used so behavior is unchanged for
//     unconfigured installations.
//   - Roles 3 (Admin) and 4 (Super Admin) default to unrestricted (all caps),
//     so an admin-level user is never accidentally blocked if role_capabilities
//     is absent from the DB.
package permission

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/gogf/gf/v2/os/gcache"

	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/service"
)

// ── default capabilities ──────────────────────────────────────────────────────
// Mirrors DEFAULT_ROLE_CAPABILITIES in frontend/app/config/permissions.ts.
// nil means "unrestricted" (all capabilities granted by default).

var defaultCaps = map[int][]string{
	1: {}, // Subscriber — no capabilities
	2: {   // Editor — content management (own resources), no cross-user or system ops
		"access_admin",
		"read_private_posts", "edit_posts", "publish_posts", "delete_posts",
		"edit_pages", "publish_pages",
		"manage_categories",
		"upload_files", "manage_media",
		"moderate_comments", "manage_comments",
		"list_users",
	},
	3: nil, // Admin — unrestricted by default (customisable via Admin UI)
	4: nil, // Super Admin — always unrestricted
}

// ── cache ─────────────────────────────────────────────────────────────────────

var capCache = gcache.New()

const capCacheTTL = time.Minute
const capCacheKey = "role_capabilities"

// InvalidateCapCache removes the cached role_capabilities so the next request
// re-reads from the database.  Call this after saving role_capabilities.
func InvalidateCapCache(ctx context.Context) {
	_, _ = capCache.Remove(ctx, capCacheKey)
}

// ── implementation ────────────────────────────────────────────────────────────

type sPermission struct{}

func New() service.IPermission { return &sPermission{} }

func init() {
	service.RegisterPermission(New())
}

// Can reports whether role has the named capability.
//
// Lookup order for custom roles (ID > 4):
//  1. Own entry in role_capabilities (e.g. "5") — Admin UI sets this per custom role.
//  2. Base role entry (e.g. "2" for baseRoleId=2) — inherit from the base role config.
//  3. Hardcoded default for the base role.
//
// Built-in roles (1-4) skip step 1 and go straight to step 2 / 3.
func (s *sPermission) Can(ctx context.Context, role int, cap string) bool {
	// Step 1 — for custom roles, try their own capability entry first.
	if role > 4 {
		if caps, found := loadCaps(ctx, role); found {
			return slices.Contains(caps, cap)
		}
		// No own entry → fall through to base role lookup.
	}

	// Step 2 — resolve to base role (1-4) and check its DB entry.
	resolved := middleware.ResolveRole(ctx, role)
	caps, found := loadCaps(ctx, resolved)
	if !found {
		// Step 3 — no DB entry at all → use hardcoded defaults.
		def, ok := defaultCaps[resolved]
		if !ok {
			return false
		}
		if def == nil {
			return true // unrestricted (admin / super admin default)
		}
		return slices.Contains(def, cap)
	}

	// nil in the DB map means unrestricted (e.g. admin saved with all caps selected).
	if caps == nil {
		return true
	}
	return slices.Contains(caps, cap)
}

// loadCaps loads the capability list for a resolved role (1-4) from DB/cache.
// found=false means there is no entry in role_capabilities for this role.
func loadCaps(ctx context.Context, role int) (caps []string, found bool) {
	allCaps := readAllCaps(ctx)
	if allCaps == nil {
		return nil, false
	}
	roleKey := fmt.Sprintf("%d", role)
	caps, found = allCaps[roleKey]
	return
}

// readAllCaps returns the full role_capabilities map from cache or DB.
// Returns nil if the option is absent.
func readAllCaps(ctx context.Context) map[string][]string {
	cached, _ := capCache.Get(ctx, capCacheKey)
	if cached != nil {
		m, _ := cached.Val().(map[string][]string)
		return m
	}

	var allCaps map[string][]string
	val, err := dao.Options.Ctx(ctx).
		Where("key", capCacheKey).
		Value("value")
	if err == nil && !val.IsNil() {
		raw := val.String()
		// DB value may be double-JSON-encoded (string wrapping JSON).
		var inner string
		if json.Unmarshal([]byte(raw), &inner) == nil {
			_ = json.Unmarshal([]byte(inner), &allCaps)
		} else {
			_ = json.Unmarshal([]byte(raw), &allCaps)
		}
	}

	// Cache even a nil result to avoid repeated DB queries.
	_ = capCache.Set(ctx, capCacheKey, allCaps, capCacheTTL)
	return allCaps
}
