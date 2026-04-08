package auth

import (
	"fmt"
	"net/url"

	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// OAuthController handles the OAuth redirect/callback flow.
// These endpoints do raw HTTP 302 redirects and cannot go through
// MiddlewareHandlerResponse, so they are registered manually in cmd.go.
type OAuthController struct{}

func NewOAuth() *OAuthController { return &OAuthController{} }

// Redirect — GET /api/v1/auth/oauth/{provider}/redirect?redirect=/path
// Generates a CSRF state token and redirects the browser to the provider's auth page.
func (c *OAuthController) Redirect(r *ghttp.Request) {
	ctx := r.GetCtx()
	provider := r.GetRouter("provider").String()
	redirectAfter := r.GetQuery("redirect", "/").String()

	authURL, err := service.Auth().OAuthRedirect(ctx, provider, redirectAfter)
	if err != nil {
		frontendBase := getFrontendBase(r)
		r.Response.RedirectTo(frontendBase+"/auth/callback?error="+url.QueryEscape(err.Error()), 302)
		return
	}
	r.Response.RedirectTo(authURL, 302)
}

// Callback — GET /api/v1/auth/oauth/{provider}/callback?code=xxx&state=xxx
// Exchanges the authorization code for JWT tokens and redirects to the frontend callback page.
func (c *OAuthController) Callback(r *ghttp.Request) {
	ctx := r.GetCtx()
	provider := r.GetRouter("provider").String()
	code := r.GetQuery("code").String()
	state := r.GetQuery("state").String()
	oauthError := r.GetQuery("error").String()

	frontendBase := getFrontendBase(r)
	callbackBase := frontendBase + "/auth/callback"

	if oauthError != "" {
		r.Response.RedirectTo(callbackBase+"?error="+url.QueryEscape(oauthError), 302)
		return
	}
	if code == "" || state == "" {
		r.Response.RedirectTo(callbackBase+"?error="+url.QueryEscape("missing code or state"), 302)
		return
	}

	access, refresh, expiresIn, redirectAfter, err := service.Auth().OAuthCallback(ctx, provider, code, state)
	if err != nil {
		r.Response.RedirectTo(callbackBase+"?error="+url.QueryEscape(err.Error()), 302)
		return
	}
	if redirectAfter == "" {
		redirectAfter = "/"
	}

	q := url.Values{}
	q.Set("access_token", access)
	q.Set("refresh_token", refresh)
	q.Set("expires_in", fmt.Sprintf("%d", expiresIn))
	q.Set("redirect", redirectAfter)
	r.Response.RedirectTo(callbackBase+"?"+q.Encode(), 302)
}

func getFrontendBase(r *ghttp.Request) string {
	val, _ := g.Cfg().Get(r.GetCtx(), "auth.oauth.frontendBase")
	if s := val.String(); s != "" {
		return s
	}
	return "http://localhost:3000"
}
