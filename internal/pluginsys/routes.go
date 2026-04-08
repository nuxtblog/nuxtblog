package pluginsys

import (
	"github.com/gogf/gf/v2/net/ghttp"
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
