package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/golang-jwt/jwt/v5"
)

type jwtClaims struct {
	UserID int64  `json:"user_id"`
	Role   int    `json:"role"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func getSecret(ctx context.Context) []byte {
	val, _ := g.Cfg().Get(ctx, "auth.jwtSecret")
	s := val.String()
	if s == "" {
		s = "change-me-in-production"
	}
	return []byte(s)
}

func ParseBearerToken(r *ghttp.Request) (*jwtClaims, error) {
	header := r.GetHeader("Authorization")
	tokenStr := strings.TrimPrefix(header, "Bearer ")
	if tokenStr == "" || tokenStr == header {
		return nil, fmt.Errorf("no token")
	}
	claims := &jwtClaims{}
	_, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return getSecret(r.GetCtx()), nil
	})
	if err != nil {
		return nil, err
	}
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, fmt.Errorf("token expired")
	}
	if claims.Type != "access" {
		return nil, fmt.Errorf("invalid token type")
	}
	return claims, nil
}

// AuthRequired rejects requests without a valid JWT
func AuthRequired(r *ghttp.Request) {
	claims, err := ParseBearerToken(r)
	if err != nil {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "unauthorized: " + err.Error(),
			"data":    nil,
		})
		return
	}
	r.SetCtxVar("user_id", claims.UserID)
	r.SetCtxVar("user_role", claims.Role)
	r.Middleware.Next()
}

// AuthOptional sets user_id if a valid JWT is present, but does not reject
func AuthOptional(r *ghttp.Request) {
	if claims, err := ParseBearerToken(r); err == nil {
		r.SetCtxVar("user_id", claims.UserID)
		r.SetCtxVar("user_role", claims.Role)
	}
	r.Middleware.Next()
}

// GetCurrentUserID extracts user_id set by AuthRequired/AuthOptional
func GetCurrentUserID(ctx context.Context) (int64, bool) {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0, false
	}
	uid := r.GetCtxVar("user_id").Int64()
	return uid, uid > 0
}

// GetCurrentUserRole extracts the raw user_role set by AuthRequired/AuthOptional.
// Returns 0 if no JWT is present.
func GetCurrentUserRole(ctx context.Context) int {
	r := ghttp.RequestFromCtx(ctx)
	if r == nil {
		return 0
	}
	return r.GetCtxVar("user_role").Int()
}

// AdminWriteRequired allows all GET/HEAD requests through freely.
// For write methods (POST/PUT/DELETE/PATCH) it requires a valid JWT with role >= 2 (editor/admin).
// Specific paths are whitelisted as public writes (comment creation, view increment, etc.).
func AdminWriteRequired(r *ghttp.Request) {
	method := r.Method
	// reads are always public
	if method == "GET" || method == "HEAD" || method == "OPTIONS" {
		r.Middleware.Next()
		return
	}

	path := r.URL.Path
	// whitelisted public write endpoints
	switch {
	case path == "/api/v1/auth/login",
		path == "/api/v1/auth/register",
		path == "/api/v1/auth/refresh",
		path == "/api/v1/auth/logout",
		path == "/api/v1/comments" && method == "POST", // guest comment creation
		strings.HasSuffix(path, "/view"),
		strings.HasSuffix(path, "/verify-password"):
		r.Middleware.Next()
		return
	}

	// authenticated-user-only endpoints (any logged-in role)
	uid := r.GetCtxVar("user_id").Int64()
	role := r.GetCtxVar("user_role").Int()
	switch {
	case path == "/api/v1/medias/upload",                         // any user can upload media
		strings.HasPrefix(path, "/api/v1/users/") && method == "PUT", // users can update own profile
		strings.HasPrefix(path, "/api/v1/reactions/"),                 // likes / bookmarks
		strings.HasPrefix(path, "/api/v1/checkin"):                    // checkin
		if uid == 0 {
			r.Response.WriteJsonExit(g.Map{
				"code":    401,
				"message": "unauthorized",
				"data":    nil,
			})
			return
		}
		r.Middleware.Next()
		return
	}

	// require editor or admin role for everything else
	// ResolveRole maps custom role IDs (>4) to their base system role so that
	// a custom role with baseRoleId=1 is not accidentally treated as Editor+.
	resolvedRole := ResolveRole(r.GetCtx(), role)
	if uid == 0 || resolvedRole < RoleEditor {
		r.Response.WriteJsonExit(g.Map{
			"code":    401,
			"message": "unauthorized",
			"data":    nil,
		})
		return
	}
	r.Middleware.Next()
}
