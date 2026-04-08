package verifycode

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	v1 "github.com/nuxtblog/nuxtblog/api/verifycode/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/notify"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"

	"github.com/gogf/gf/v2/frame/g"
)

type sVerifycode struct{}

func New() service.IVerifycode { return &sVerifycode{} }

func init() {
	service.RegisterVerifycode(New())
}

// Send — dispatches a 6-digit registration code to the target via email or SMS.
func (s *sVerifycode) Send(ctx context.Context, req *v1.VerifycodeSendReq) (*v1.VerifycodeSendRes, error) {
	if s.GetRegisterVerifyMode(ctx) != req.Type {
		return nil, errors.New(g.I18n().T(ctx, "verifycode.mode_disabled"))
	}

	// Rate limit: one code per target per 60 seconds
	recent, _ := dao.VerificationCodes.Ctx(ctx).
		Where("target", req.Target).
		Where("type", req.Type+"_register").
		Where("created_at > ?", time.Now().Add(-60*time.Second).Format("2006-01-02 15:04:05")).
		Count()
	if recent > 0 {
		return nil, errors.New(g.I18n().T(ctx, "verifycode.send_too_frequent"))
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	expiresAt := time.Now().Add(10 * time.Minute).Format("2006-01-02 15:04:05")

	if _, err := dao.VerificationCodes.Ctx(ctx).Data(g.Map{
		"id":         idgen.New(),
		"target":     req.Target,
		"code":       code,
		"type":       req.Type + "_register",
		"expires_at": expiresAt,
	}).Insert(); err != nil {
		return nil, errors.New(g.I18n().T(ctx, "verifycode.generate_failed"))
	}

	siteName := s.getSiteName(ctx)

	switch req.Type {
	case "email":
		body := fmt.Sprintf(`<!DOCTYPE html><html><body style="font-family:sans-serif;max-width:500px;margin:0 auto;padding:24px">
<h2 style="color:#7c3aed">%s</h2>
<p>%s</p>
<div style="font-size:36px;font-weight:700;letter-spacing:8px;color:#7c3aed;margin:16px 0">%s</div>
<p style="color:#666">%s</p>
<hr style="margin:24px 0;border:none;border-top:1px solid #eee">
<p style="color:#999;font-size:12px">%s</p>
</body></html>`,
			g.I18n().Tf(ctx, "verifycode.email_body_title", siteName),
			g.I18n().T(ctx, "verifycode.email_body_intro"),
			code,
			g.I18n().T(ctx, "verifycode.email_body_expiry"),
			g.I18n().Tf(ctx, "verifycode.email_body_auto", siteName),
		)
		if err := notify.SendEmail(ctx, req.Target, g.I18n().Tf(ctx, "verifycode.email_subject", siteName), body); err != nil {
			return nil, fmt.Errorf("%s: %w", g.I18n().T(ctx, "verifycode.email_send_failed"), err)
		}
	case "sms":
		content := g.I18n().Tf(ctx, "verifycode.sms_content", siteName, code)
		if err := notify.SendSMS(ctx, req.Target, content); err != nil {
			return nil, fmt.Errorf("%s: %w", g.I18n().T(ctx, "verifycode.sms_send_failed"), err)
		}
	}

	return &v1.VerifycodeSendRes{ExpiresIn: 600}, nil
}

// CheckCode validates a code submitted during registration.
func (s *sVerifycode) CheckCode(ctx context.Context, target, codeType, code string) error {
	now := time.Now().Format("2006-01-02 15:04:05")
	type row struct {
		Id   int64  `orm:"id"`
		Code string `orm:"code"`
	}
	var r row
	err := dao.VerificationCodes.Ctx(ctx).
		Where("target", target).
		Where("type", codeType+"_register").
		WhereNull("used_at").
		Where("expires_at > ?", now).
		OrderDesc("id").
		Limit(1).
		Scan(&r)
	if err != nil || r.Id == 0 {
		return errors.New(g.I18n().T(ctx, "verifycode.invalid_or_expired"))
	}
	if r.Code != code {
		return errors.New(g.I18n().T(ctx, "verifycode.wrong_code"))
	}
	_, _ = dao.VerificationCodes.Ctx(ctx).
		Where("id", r.Id).
		Update(g.Map{"used_at": now})
	return nil
}

// GetRegisterVerifyMode reads auth_register_verify.mode from options ("none"|"email"|"sms").
func (s *sVerifycode) GetRegisterVerifyMode(ctx context.Context) string {
	val, err := dao.Options.Ctx(ctx).Where("key", "auth_register_verify").Value("value")
	if err != nil || val.IsEmpty() {
		return "none"
	}
	var cfg struct{ Mode string `json:"mode"` }
	if err := json.Unmarshal([]byte(val.String()), &cfg); err != nil {
		return "none"
	}
	return cfg.Mode
}

func (s *sVerifycode) getSiteName(ctx context.Context) string {
	val, err := dao.Options.Ctx(ctx).Where("key", "notify_email").Value("value")
	if err != nil || val.IsEmpty() {
		return "Blog"
	}
	var cfg struct{ SiteName string `json:"site_name"` }
	if err := json.Unmarshal([]byte(val.String()), &cfg); err != nil || cfg.SiteName == "" {
		return "Blog"
	}
	return cfg.SiteName
}
