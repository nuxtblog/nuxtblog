package middleware

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
)

// ── Custom role level cache ───────────────────────────────────────────────────

// customRoleCache holds parsed custom_role_defs from the options table.
// TTL is 1 minute; invalidated on any write that changes role definitions.
var customRoleCache = gcache.New()

const customRoleCacheTTL = time.Minute

type customRoleDef struct {
	ID         int `json:"id"`
	BaseRoleID int `json:"baseRoleId"`
}

// ResolveRole maps a raw role integer to the effective system role level (1-4).
// Built-in roles 1-4 are returned as-is.
// Custom roles (ID > 4) look up their baseRoleId from the options table
// (result is cached for customRoleCacheTTL).
// Falls back to RoleSubscriber if the role is unknown.
func ResolveRole(ctx context.Context, role int) int {
	if role >= 1 && role <= 4 {
		return role
	}

	// Try cache first.
	const cacheKey = "custom_role_defs"
	cached, _ := customRoleCache.Get(ctx, cacheKey)

	var defs []customRoleDef
	if cached != nil {
		defs, _ = cached.Val().([]customRoleDef)
	} else {
		// Load from options table.
		val, err := g.DB().Model("options").Ctx(ctx).
			Where("key", "custom_role_defs").
			Value("value")
		if err == nil && !val.IsNil() {
			_ = json.Unmarshal([]byte(val.String()), &defs)
		}
		_ = customRoleCache.Set(ctx, cacheKey, defs, customRoleCacheTTL)
	}

	for _, d := range defs {
		if d.ID == role {
			if d.BaseRoleID >= 1 && d.BaseRoleID <= 4 {
				return d.BaseRoleID
			}
		}
	}
	// Unknown custom role — treat as least-privileged.
	return RoleSubscriber
}

// InvalidateCustomRoleCache removes the cached role definitions so the next
// request re-reads from the database. Call this after saving custom_role_defs.
func InvalidateCustomRoleCache(ctx context.Context) {
	_, _ = customRoleCache.Remove(ctx, "custom_role_defs")
}



// ── Role constants ─────────────────────────────────────────────────────────────

const (
	RoleSubscriber = 1
	RoleEditor     = 2
	RoleAdmin      = 3
	RoleSuperAdmin = 4
)

// ── Route rule ────────────────────────────────────────────────────────────────

// routeRule declares the minimum role required for a matched route.
//
// Pattern rules:
//   - Exact segment match:  /api/v1/users
//   - Single wildcard (*):  /api/v1/users/*   matches /api/v1/users/123
//   - Any-suffix wildcard (**): /api/v1/posts/**  matches any sub-path
//
// Method: HTTP verb or "*" to match any method.
// MinRole: minimum role integer required (1-4). 0 = skip.
type routeRule struct {
	Method  string
	Pattern string
	MinRole int
}

// routeRules is evaluated top-to-bottom; the first matching rule wins.
// Rules that require admin(3) should come before broader editor(2) rules
// when the paths overlap.
var routeRules = []routeRule{

	// ── Site options — admin only ──────────────────────────────────────────────
	// PUT  /api/v1/options/{key}   → set option value
	// DELETE /api/v1/options/{key} → delete option
	{Method: "PUT",    Pattern: "/api/v1/options/*",     MinRole: RoleAdmin},
	{Method: "DELETE", Pattern: "/api/v1/options/*",     MinRole: RoleAdmin},

	// ── Navigation menus — admin only ─────────────────────────────────────────
	{Method: "POST",   Pattern: "/api/v1/nav-menus",     MinRole: RoleAdmin},
	{Method: "POST",   Pattern: "/api/v1/nav-menus/*",   MinRole: RoleAdmin},
	{Method: "PUT",    Pattern: "/api/v1/nav-menus/*",   MinRole: RoleAdmin},
	{Method: "DELETE", Pattern: "/api/v1/nav-menus/*",   MinRole: RoleAdmin},

	// ── Users — list / create / delete require admin ───────────────────────────
	// Note: PUT /users/{id} is intentionally excluded here because users
	// can update their own profile (whitelisted in AdminWriteRequired).
	{Method: "GET",    Pattern: "/api/v1/users",         MinRole: RoleAdmin},
	{Method: "POST",   Pattern: "/api/v1/users",         MinRole: RoleAdmin},
	{Method: "DELETE", Pattern: "/api/v1/users/*",       MinRole: RoleAdmin},

	// ── Taxonomy — editor+ ────────────────────────────────────────────────────
	{Method: "POST",   Pattern: "/api/v1/terms",         MinRole: RoleEditor},
	{Method: "PUT",    Pattern: "/api/v1/terms/*",       MinRole: RoleEditor},
	{Method: "DELETE", Pattern: "/api/v1/terms/*",       MinRole: RoleEditor},
	{Method: "POST",   Pattern: "/api/v1/taxonomies",    MinRole: RoleEditor},
	{Method: "PUT",    Pattern: "/api/v1/taxonomies/*",  MinRole: RoleEditor},
	{Method: "DELETE", Pattern: "/api/v1/taxonomies/*",  MinRole: RoleEditor},

	// ── Posts — editor+ ───────────────────────────────────────────────────────
	{Method: "POST",   Pattern: "/api/v1/posts",          MinRole: RoleEditor},
	{Method: "PUT",    Pattern: "/api/v1/posts/*",        MinRole: RoleEditor},
	{Method: "DELETE", Pattern: "/api/v1/posts/*",        MinRole: RoleEditor},
	{Method: "POST",   Pattern: "/api/v1/posts/batch",    MinRole: RoleEditor},

	// ── Comments — moderation requires editor+ ────────────────────────────────
	{Method: "PUT",    Pattern: "/api/v1/comments/*",    MinRole: RoleEditor},
	{Method: "DELETE", Pattern: "/api/v1/comments/*",    MinRole: RoleEditor},

	// ── Commerce admin — admin only ──────────────────────────────────────────
	{Method: "GET",    Pattern: "/admin/orders",                   MinRole: RoleAdmin},
	{Method: "POST",   Pattern: "/admin/orders/*/refund",          MinRole: RoleAdmin},
	{Method: "GET",    Pattern: "/admin/revenue/stats",            MinRole: RoleAdmin},
	// Plugin admin routes (wallet, credits, membership)
	{Method: "*",      Pattern: "/api/plugin/wallet/admin/*",      MinRole: RoleAdmin},
	{Method: "*",      Pattern: "/api/plugin/credits/admin/*",     MinRole: RoleAdmin},
	{Method: "*",      Pattern: "/api/plugin/membership/admin/*",  MinRole: RoleAdmin},

	// ── Media library — GET requires editor+ (media library is admin-facing) ──
	// POST /medias/upload is whitelisted in AdminWriteRequired for any auth user.
	{Method: "GET",    Pattern: "/api/v1/medias",        MinRole: RoleEditor},
	{Method: "GET",    Pattern: "/api/v1/medias/stats",  MinRole: RoleEditor},
	{Method: "GET",    Pattern: "/api/v1/medias/*",      MinRole: RoleEditor},
	{Method: "PUT",    Pattern: "/api/v1/medias/*",      MinRole: RoleEditor},
	{Method: "DELETE", Pattern: "/api/v1/medias/*",      MinRole: RoleEditor},
}

// ── Middleware ────────────────────────────────────────────────────────────────

// RouteRBACCheck enforces per-route minimum role requirements.
// It must be mounted AFTER AuthOptional so that user_role is available in ctx.
// Requests that do not match any rule are passed through unchanged.
func RouteRBACCheck(r *ghttp.Request) {
	path   := r.URL.Path
	method := r.Method

	for _, rule := range routeRules {
		if !matchPattern(rule.Pattern, path) {
			continue
		}
		if rule.Method != "*" && rule.Method != method {
			continue
		}

		// Rule matched — check role (custom role IDs are resolved to their base level).
		role := ResolveRole(r.GetCtx(), r.GetCtxVar("user_role").Int())
		if role < rule.MinRole {
			code := 403
			msg  := "forbidden: insufficient role"
			if role == 0 {
				code = 401
				msg  = "unauthorized"
			}
			r.Response.WriteJsonExit(g.Map{
				"code":    code,
				"message": msg,
				"data":    nil,
			})
			return
		}
		break // first match wins; no need to check further rules
	}

	r.Middleware.Next()
}

// ── Pattern matching ──────────────────────────────────────────────────────────

// matchPattern reports whether the URL path matches a pattern.
//
//   *  — matches exactly one path segment (e.g. "123" or "abc")
//   ** — matches any remaining path (must be the last segment in the pattern)
func matchPattern(pattern, path string) bool {
	pSegs := splitPath(pattern)
	rSegs := splitPath(path)

	for i, p := range pSegs {
		if p == "**" {
			// matches everything remaining
			return true
		}
		if i >= len(rSegs) {
			return false
		}
		if p != "*" && p != rSegs[i] {
			return false
		}
	}
	return len(pSegs) == len(rSegs)
}

func splitPath(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

// ── Ownership check ───────────────────────────────────────────────────────────

// ownershipRule declares that for a matched route, an editor-level user (role==2)
// must own the target resource. Admin and above always bypass this check.
//
// Table:    DB table to query.
// OwnerCol: column that holds the owner's user ID.
// The resource ID is extracted from the path segment that matches "*" in Pattern.
type ownershipRule struct {
	Method   string
	Pattern  string
	Table    string
	OwnerCol string
}

// ownershipRules covers mutable operations on individually-owned resources.
// Comment moderation is intentionally excluded — editors can moderate any comment.
var ownershipRules = []ownershipRule{
	{Method: "PUT",    Pattern: "/api/v1/posts/*",   Table: "posts",   OwnerCol: "author_id"},
	{Method: "DELETE", Pattern: "/api/v1/posts/*",   Table: "posts",   OwnerCol: "author_id"},
	{Method: "PUT",    Pattern: "/api/v1/medias/*",  Table: "medias",  OwnerCol: "uploader_id"},
	{Method: "DELETE", Pattern: "/api/v1/medias/*",  Table: "medias",  OwnerCol: "uploader_id"},
}

// OwnershipCheck enforces that editor-level users (role==2) can only mutate
// resources they own. Admin (3) and super-admin (4) bypass this check entirely.
// Must be mounted AFTER RouteRBACCheck so that unauthenticated/low-role requests
// are already rejected before we hit the database.
func OwnershipCheck(r *ghttp.Request) {
	role := ResolveRole(r.GetCtx(), r.GetCtxVar("user_role").Int())

	// Admin and above own everything — skip.
	if role >= RoleAdmin {
		r.Middleware.Next()
		return
	}

	path   := r.URL.Path
	method := r.Method

	for _, rule := range ownershipRules {
		if rule.Method != "*" && rule.Method != method {
			continue
		}
		if !matchPattern(rule.Pattern, path) {
			continue
		}

		// Extract the resource ID from the wildcard segment.
		resourceID, ok := extractWildcard(rule.Pattern, path)
		if !ok || resourceID == "" {
			break
		}

		userID := r.GetCtxVar("user_id").Int64()
		ctx    := r.GetCtx()

		count, err := g.DB().Model(rule.Table).Ctx(ctx).
			Where("id", resourceID).
			Where(rule.OwnerCol, userID).
			WhereNull("deleted_at").
			Count()

		if err != nil || count == 0 {
			r.Response.WriteJsonExit(g.Map{
				"code":    403,
				"message": "forbidden: you do not own this resource",
				"data":    nil,
			})
			return
		}
		break // first match wins
	}

	r.Middleware.Next()
}

// extractWildcard returns the path segment that corresponds to the first "*"
// wildcard in the pattern.
func extractWildcard(pattern, path string) (string, bool) {
	pSegs := splitPath(pattern)
	rSegs := splitPath(path)
	for i, p := range pSegs {
		if p == "*" {
			if i < len(rSegs) {
				return rSegs[i], true
			}
			return "", false
		}
	}
	return "", false
}
