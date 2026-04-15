package logic

import (
	_ "github.com/nuxtblog/nuxtblog/internal/listener" // register all event listeners
	_ "github.com/nuxtblog/nuxtblog/internal/notify"   // register email/SMS notification channels

	_ "github.com/nuxtblog/nuxtblog/internal/logic/auth"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/history"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/follow"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/token"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/verifycode"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/comment"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/media"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/announcement"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/notification"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/option"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/permission"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/post"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/reaction"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/report"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/taxonomy"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/user"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/message"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/plugin"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/doc"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/moment"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/ai"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/payment"

	_ "github.com/nuxtblog/nuxtblog/internal/logic/wallet"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/credits"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/order"
	_ "github.com/nuxtblog/nuxtblog/internal/logic/membership"
)
