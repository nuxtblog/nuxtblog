package notify

import (
	"context"
	"encoding/json"
	"fmt"
	"net/smtp"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/nuxtblog/nuxtblog/internal/consts"
)

func init() {
	Register(&EmailChannel{})
}

// EmailConfig is stored in the options table under key "notify_email".
type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	SiteName string `json:"site_name"`
	SiteURL  string `json:"site_url"`
}

type EmailChannel struct{}

func (c *EmailChannel) Name() string { return consts.NotifyChannelEmail }

func (c *EmailChannel) Enabled(ctx context.Context) bool {
	cfg := loadEmailConfig(ctx)
	return cfg != nil && cfg.Host != ""
}

func (c *EmailChannel) Send(ctx context.Context, msg Message) error {
	cfg := loadEmailConfig(ctx)
	if cfg == nil || cfg.Host == "" {
		return nil
	}

	port := cfg.Port
	if port == 0 {
		port = 587
	}
	site := cfg.SiteName
	if site == "" {
		site = "Blog"
	}

	// Fetch recipient email
	type userRow struct{ Email string `orm:"email"` }
	var u userRow
	if err := g.DB().Ctx(ctx).Model("users").Where("id", msg.UserID).Fields("email").Scan(&u); err != nil || u.Email == "" {
		return nil
	}

	subject, body := formatEmail(msg, site, cfg.SiteURL)
	addr := fmt.Sprintf("%s:%d", cfg.Host, port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	raw := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s",
		cfg.From, u.Email, subject, body,
	)
	return smtp.SendMail(addr, auth, cfg.From, []string{u.Email}, []byte(raw))
}

// loadEmailConfig reads EmailConfig from options table.
func loadEmailConfig(ctx context.Context) *EmailConfig {
	val, err := g.DB().Ctx(ctx).Model("options").Where("key", "notify_email").Value("value")
	if err != nil || val.IsEmpty() {
		return nil
	}
	var cfg EmailConfig
	if err := json.Unmarshal([]byte(val.String()), &cfg); err != nil {
		return nil
	}
	return &cfg
}

func formatEmail(msg Message, site, baseURL string) (subject, body string) {
	actor := msg.ActorName
	if actor == "" {
		actor = "有人"
	}
	switch msg.Type {
	case "follow":
		subject = fmt.Sprintf("[%s] %s 开始关注了你", site, actor)
		link := baseURL
		if msg.ActorID != nil {
			link = fmt.Sprintf("%s/user/%d", baseURL, *msg.ActorID)
		}
		body = htmlWrap(fmt.Sprintf(`<p><strong>%s</strong> 开始关注了你。</p>
<p><a href="%s" style="color:#7c3aed">查看 Ta 的主页 →</a></p>`, actor, link), site)
	case "like":
		subject = fmt.Sprintf("[%s] %s 点赞了你的文章", site, actor)
		body = htmlWrap(fmt.Sprintf(`<p><strong>%s</strong> 点赞了你的文章
<a href="%s%s" style="color:#7c3aed">《%s》</a></p>`, actor, baseURL, msg.ObjectLink, msg.ObjectTitle), site)
	case "comment":
		subject = fmt.Sprintf("[%s] %s 评论了你的文章", site, actor)
		body = htmlWrap(fmt.Sprintf(`<p><strong>%s</strong> 评论了你的文章
<a href="%s%s" style="color:#7c3aed">《%s》</a>：</p>
<blockquote style="border-left:3px solid #7c3aed;padding-left:12px;color:#555">%s</blockquote>`,
			actor, baseURL, msg.ObjectLink, msg.ObjectTitle, msg.Content), site)
	case "reply":
		subject = fmt.Sprintf("[%s] %s 回复了你的评论", site, actor)
		body = htmlWrap(fmt.Sprintf(`<p><strong>%s</strong> 回复了你的评论：</p>
<blockquote style="border-left:3px solid #7c3aed;padding-left:12px;color:#555">%s</blockquote>`,
			actor, msg.Content), site)
	case "system":
		subject = fmt.Sprintf("[%s] 系统通知", site)
		body = htmlWrap(fmt.Sprintf(`<p>%s</p>`, msg.Content), site)
	default:
		subject = fmt.Sprintf("[%s] 你有一条新通知", site)
		body = htmlWrap(fmt.Sprintf(`<p>%s %s</p>`, actor, msg.Content), site)
	}
	return
}

func htmlWrap(content, site string) string {
	return fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:sans-serif;max-width:600px;margin:0 auto;padding:24px">
<h2 style="color:#7c3aed;margin-bottom:16px">%s</h2>
%s
<hr style="margin:24px 0;border:none;border-top:1px solid #eee">
<p style="color:#999;font-size:12px">此邮件由 %s 自动发送，请勿回复。</p>
</body></html>`, site, content, site)
}
