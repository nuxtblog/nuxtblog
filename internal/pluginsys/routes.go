package pluginsys

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/sdk"
)

// ─── Plugin HTTP Route Types (Phase 2.7) ────────────────────────────────────

// PluginRequest is the object passed to a plugin route handler.
type PluginRequest struct {
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Query    map[string]string `json:"query"`
	Body     any               `json:"body"`
	Headers  map[string]string `json:"headers"`
	UserID   int               `json:"userId,omitempty"`
	UserRole int               `json:"userRole,omitempty"`
}

// PluginResponse is the object returned by a plugin route handler.
type PluginResponse struct {
	Status  int               `json:"status"`
	Body    any               `json:"body"`
	Headers map[string]string `json:"headers,omitempty"`
}

// getUserIDFromRequest extracts the user ID from the request context.
// Returns 0 if not authenticated.
func getUserIDFromRequest(r *ghttp.Request) int {
	uid := r.GetCtxVar("user_id")
	if !uid.IsNil() && !uid.IsEmpty() {
		return uid.Int()
	}
	return 0
}

// getUserRoleFromRequest extracts the user role from the request context.
// Returns 0 if not authenticated. Role values: 1=subscriber, 2=editor, 3=admin.
func getUserRoleFromRequest(r *ghttp.Request) int {
	role := r.GetCtxVar("user_role")
	if !role.IsNil() && !role.IsEmpty() {
		return role.Int()
	}
	return 0
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

	// Allow /api/ prefix to pass through; auto-prefix others under /api/plugin/{id}
	if !strings.HasPrefix(path, "/api/") {
		path = fmt.Sprintf("/api/plugin/%s%s", reg.pluginID, path)
	}

	reg.entries = append(reg.entries, routeEntry{
		method:  strings.ToUpper(method),
		path:    path,
		handler: handler,
		auth:    auth,
	})
}
