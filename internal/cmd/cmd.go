package cmd

import (
	"context"
	"net/http"
	"os"

	_ "github.com/nuxtblog/nuxtblog/internal/logic"
	_ "github.com/nuxtblog/nuxtblog/builtin" // static Go-native plugin registration
	"github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/text/gstr"

	aiCtrl "github.com/nuxtblog/nuxtblog/internal/controller/ai"
	orderCtrl "github.com/nuxtblog/nuxtblog/internal/controller/order"
	paymentCtrl "github.com/nuxtblog/nuxtblog/internal/controller/payment"
	walletCtrl "github.com/nuxtblog/nuxtblog/internal/controller/wallet"
	announcementCtrl "github.com/nuxtblog/nuxtblog/internal/controller/announcement"
	"github.com/nuxtblog/nuxtblog/internal/controller/auth"
	docCtrl "github.com/nuxtblog/nuxtblog/internal/controller/doc"
	momentCtrl "github.com/nuxtblog/nuxtblog/internal/controller/moment"
	"github.com/nuxtblog/nuxtblog/internal/controller/checkin"
	"github.com/nuxtblog/nuxtblog/internal/controller/comment"
	"github.com/nuxtblog/nuxtblog/internal/controller/follow"
	historyCtrl "github.com/nuxtblog/nuxtblog/internal/controller/history"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/controller/media"
	"github.com/nuxtblog/nuxtblog/internal/controller/notification"
	"github.com/nuxtblog/nuxtblog/internal/controller/option"
	"github.com/nuxtblog/nuxtblog/internal/controller/post"
	"github.com/nuxtblog/nuxtblog/internal/controller/reaction"
	messageCtrl "github.com/nuxtblog/nuxtblog/internal/controller/message"
	pluginCtrl "github.com/nuxtblog/nuxtblog/internal/controller/plugin"
	reportCtrl "github.com/nuxtblog/nuxtblog/internal/controller/report"
	"github.com/nuxtblog/nuxtblog/internal/controller/site"
	"github.com/nuxtblog/nuxtblog/internal/controller/taxonomy"
	"github.com/nuxtblog/nuxtblog/internal/controller/token"
	"github.com/nuxtblog/nuxtblog/internal/controller/user"
	"github.com/nuxtblog/nuxtblog/internal/controller/verifycode"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			middleware.InitLogger()

			if err = autoMigrate(ctx); err != nil {
				g.Log().Warningf(ctx, "autoMigrate warning: %v", err)
			}

			// Plugin manager — handles builtin (Go) + external (Goja JS) plugins
			pluginMgr := pluginsys.New()
			eng.RegisterGoNativeManager(pluginMgr)
			// Layer 2: load compiled Go plugins (builtin)
			if err := pluginMgr.LoadStatic(ctx); err != nil {
				g.Log().Warningf(ctx, "[pluginmgr] LoadStatic: %v", err)
			}
			// Layer 1+: load external JS/YAML/UI plugins from data/plugins/
			if err := pluginMgr.LoadExternal(ctx, "data/plugins"); err != nil {
				g.Log().Warningf(ctx, "[pluginmgr] LoadExternal: %v", err)
			}

			// Layer 0: load declarative YAML plugins from data/plugins/yaml/
			eng.LoadYAMLPlugins(ctx, "data/plugins/yaml")

			// Ensure plugin migration table exists (used by all layers)
			eng.EnsureMigrationTable(ctx)
			// Register event bus listeners that fan out to all plugin layers
			eng.RegisterEventListeners()

			// Cron: auto-publish scheduled posts every minute (GoFrame gcron uses 6-field format: sec min hour day month week)
			if _, err = gcron.Add(ctx, "0 * * * * *", func(ctx context.Context) {
				service.Post().PublishScheduled(ctx)
			}, "publish-scheduled"); err != nil {
				g.Log().Warningf(ctx, "register publish-scheduled cron: %v", err)
			}

			// Cron: check expired memberships every hour
			if _, err = gcron.Add(ctx, "0 0 * * * *", func(ctx context.Context) {
				service.Membership().ExpireCheck(ctx)
			}, "membership-expire-check"); err != nil {
				g.Log().Warningf(ctx, "register membership-expire-check cron: %v", err)
			}

			s := g.Server()

			// Plugin routes (builtin + JS)
			pluginMgr.RegisterRoutes(s)
			pluginMgr.RegisterAssetRoutes(s, "data/plugins")
			eng.RegisterServiceProxies(s)

			// Serve locally uploaded files at /uploads (only when driver=local)
			storageDriver, _ := g.Cfg().Get(ctx, "storage.driver")
			if storageDriver.String() == "" || storageDriver.String() == "local" {
				uploadPath, _ := g.Cfg().Get(ctx, "storage.local.uploadPath")
				p := uploadPath.String()
				if p == "" {
					p = "./uploads"
				}
				if err = os.MkdirAll(p, 0755); err != nil {
					g.Log().Warningf(ctx, "create upload dir: %v", err)
				}
				s.AddStaticPath("/uploads", p)
			}

			corsOrigins, _ := g.Cfg().Get(ctx, "server.corsOrigins")
			allowedOrigin := corsOrigins.String()
			if allowedOrigin == "" {
				allowedOrigin = "*"
			}
			// Colorful access logger — global, covers all routes
			s.Use(middleware.AccessLogger)

			// Global CORS — configurable via server.corsOrigins (comma-separated or *)
			s.Use(func(r *ghttp.Request) {
				origin := r.GetHeader("Origin")
				if allowedOrigin == "*" {
					r.Response.Header().Set("Access-Control-Allow-Origin", "*")
				} else {
					for _, o := range gstr.SplitAndTrim(allowedOrigin, ",") {
						if o == origin {
							r.Response.Header().Set("Access-Control-Allow-Origin", origin)
							r.Response.Header().Set("Vary", "Origin")
							break
						}
					}
				}
				r.Response.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
				r.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With")
				if r.Method == http.MethodOptions {
					r.Response.WriteHeader(http.StatusNoContent)
					r.Exit()
				}
				r.Middleware.Next()
			})

			// OAuth redirect/callback: raw HTTP 302, must be outside MiddlewareHandlerResponse.
			// /providers is a normal JSON endpoint auto-bound via group.Bind(auth.NewV1()).
			oauthCtrl := auth.NewOAuth()
			s.Group("/api/v1/auth/oauth", func(og *ghttp.RouterGroup) {
				og.GET("/{provider}/redirect", oauthCtrl.Redirect)
				og.GET("/{provider}/callback", oauthCtrl.Callback)
			})

			s.Group("/api/v1", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.LangMiddleware)
				group.Middleware(ghttp.MiddlewareHandlerResponse)
				group.Middleware(middleware.AuthOptional)
				registerPublicRoutes(group)
				registerAdminRoutes(group)
				registerAuthenticatedRoutes(group)
			})
			s.Run()
			return nil
		},
	}
)

// registerPublicRoutes binds auth endpoints and public read-only routes.
func registerPublicRoutes(group *ghttp.RouterGroup) {
	group.Bind(
		auth.NewV1(),
		verifycode.NewV1(),
		follow.NewPublicV1(),
		site.NewV1(),
		pluginCtrl.NewPublic(),
	)
}

// registerAdminRoutes binds write endpoints that require editor/admin role.
func registerAdminRoutes(group *ghttp.RouterGroup) {
	group.Group("/", func(g *ghttp.RouterGroup) {
		g.Middleware(middleware.AdminWriteRequired)
		g.Middleware(middleware.RouteRBACCheck)
		g.Middleware(middleware.OwnershipCheck)
		g.Bind(
			user.NewV1(),
			post.NewV1(),
			taxonomy.NewV1(),
			comment.NewV1(),
			media.NewV1(),
			option.NewV1(),
			pluginCtrl.New(),
			announcementCtrl.NewAdmin(),
			docCtrl.NewV1(),
			aiCtrl.New(),
			paymentCtrl.New(),
			orderCtrl.New(),
			walletCtrl.New(),
		)
	})
}

// registerAuthenticatedRoutes binds endpoints that require any valid JWT.
func registerAuthenticatedRoutes(group *ghttp.RouterGroup) {
	group.Group("/", func(protected *ghttp.RouterGroup) {
		protected.Middleware(middleware.AuthRequired)
		protected.Bind(
			reaction.NewV1(),
			notification.NewV1(),
			checkin.NewV1(),
			token.NewV1(),
			follow.NewAuthV1(),
			reportCtrl.New(),
			historyCtrl.New(),
			messageCtrl.New(),
			announcementCtrl.NewPublic(),
			momentCtrl.NewV1(),
		)
	})
}
