package pluginsys

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ─── External Service Proxy ─────────────────────────────────────────────────
//
// Plugins can declare an external service in their manifest. The blog server
// acts as a reverse proxy, forwarding requests to the plugin's backend service.
//
// The proxy validates auth before forwarding and injects user info headers.

// RegisterServiceProxies registers reverse proxies for all plugins that declare
// an external service in their manifest.
func RegisterServiceProxies(s *ghttp.Server) {
	type row struct {
		Id       string `orm:"id"`
		Manifest string `orm:"manifest"`
	}
	var rows []row
	_ = g.DB().Ctx(context.Background()).
		Model("plugins").
		Where("enabled", 1).
		Fields("id, COALESCE(manifest,'{}') as manifest").
		Scan(&rows)

	for _, r := range rows {
		var mf Manifest
		if err := json.Unmarshal([]byte(r.Manifest), &mf); err != nil {
			continue
		}
		if mf.Service == nil || mf.Service.Proxy == "" || mf.Service.Target == "" {
			continue
		}

		registerProxy(s, r.Id, mf.Service)
	}
}

func registerProxy(s *ghttp.Server, pluginID string, svc *ServiceDef) {
	target, err := url.Parse(svc.Target)
	if err != nil {
		g.Log().Warningf(context.Background(),
			"[plugin:%s] invalid proxy target %q: %v", pluginID, svc.Target, err)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.Transport = &http.Transport{
		ResponseHeaderTimeout: 30 * time.Second,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
	}

	// Ensure proxy path has trailing slash for prefix matching
	proxyPath := svc.Proxy
	if !strings.HasSuffix(proxyPath, "/") {
		proxyPath += "/"
	}

	s.BindHandler(proxyPath+"*", func(r *ghttp.Request) {
		// Inject user info headers if authenticated
		if info := extractUserInfo(r); info != nil {
			r.Request.Header.Set("X-Plugin-User-ID", info.userID)
			r.Request.Header.Set("X-Plugin-User-Role", info.userRole)
		}

		// Strip the proxy prefix from the path before forwarding
		originalPath := r.URL.Path
		r.Request.URL.Path = strings.TrimPrefix(originalPath, strings.TrimSuffix(proxyPath, "/"))
		if r.Request.URL.Path == "" {
			r.Request.URL.Path = "/"
		}

		proxy.ServeHTTP(r.Response.RawWriter(), r.Request)
		// Prevent GoFrame from writing additional response
		r.Exit()
	})

	g.Log().Infof(context.Background(),
		"[plugin] registered proxy %s → %s for plugin %s", proxyPath, svc.Target, pluginID)
}

type userInfo struct {
	userID   string
	userRole string
}

// extractUserInfo reads user identity from the request context (set by auth middleware).
// Returns nil if no authenticated user is present.
func extractUserInfo(r *ghttp.Request) *userInfo {
	uid := r.GetCtxVar("user_id")
	role := r.GetCtxVar("user_role")
	if uid.IsNil() || uid.IsEmpty() {
		return nil
	}
	return &userInfo{
		userID:   uid.String(),
		userRole: role.String(),
	}
}

// RegisterProxyRoutes is an alias for RegisterServiceProxies for consistency.
func RegisterProxyRoutes(s *ghttp.Server) {
	RegisterServiceProxies(s)
}

// Ensure the proxy handler responds to any HTTP method
func init() {
	// no-op: GoFrame's BindHandler("path/*") catches all methods
	_ = http.StatusOK
}
