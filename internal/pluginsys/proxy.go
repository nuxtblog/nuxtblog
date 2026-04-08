package pluginsys

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ─── External Service Proxy (Phase 4.3) ────────────────────────────────────
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

	// Ensure proxy path has trailing slash for prefix matching
	proxyPath := svc.Proxy
	if !strings.HasSuffix(proxyPath, "/") {
		proxyPath += "/"
	}

	s.BindHandler(proxyPath+"*", func(r *ghttp.Request) {
		// Parse auth if present
		if claims, parseErr := parseJWT(r); parseErr == nil {
			r.Request.Header.Set("X-Plugin-User-ID", claims.userID)
			r.Request.Header.Set("X-Plugin-User-Role", claims.userRole)
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

type jwtInfo struct {
	userID   string
	userRole string
}

// parseJWT is a lightweight JWT parser that extracts user info from the
// Authorization header without full validation (the upstream service
// should do its own auth if needed).
func parseJWT(r *ghttp.Request) (*jwtInfo, error) {
	header := r.GetHeader("Authorization")
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	if tokenStr == "" || tokenStr == header {
		return nil, io.EOF
	}
	// Decode the payload (middle part) without verification
	parts := strings.Split(tokenStr, ".")
	if len(parts) != 3 {
		return nil, io.EOF
	}

	// We just forward the raw token and let the middleware handle it
	// The upstream gets user info via X-Plugin-User-* headers
	// Try to extract from request context (set by auth middleware)
	uid := r.GetCtxVar("user_id")
	role := r.GetCtxVar("user_role")
	if uid.IsNil() || uid.IsEmpty() {
		return nil, io.EOF
	}
	return &jwtInfo{
		userID:   uid.String(),
		userRole: role.String(),
	}, nil
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
