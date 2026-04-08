package notify

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/nuxtblog/nuxtblog/internal/notify/smsprovider"
)

// SendEmail sends a single transactional email (e.g. verification code).
// It reuses the SMTP config stored in options table (notify_email).
// Returns nil silently if the email channel is not configured.
func SendEmail(ctx context.Context, to, subject, htmlBody string) error {
	cfg := loadEmailConfig(ctx)
	if cfg == nil || cfg.Host == "" {
		return fmt.Errorf("邮件服务未配置")
	}
	port := cfg.Port
	if port == 0 {
		port = 587
	}
	addr := fmt.Sprintf("%s:%d", cfg.Host, port)
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	raw := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=utf-8\r\n\r\n%s",
		cfg.From, to, subject, htmlBody,
	)
	return smtp.SendMail(addr, auth, cfg.From, []string{to}, []byte(raw))
}

// SendSMS sends a transactional SMS message.
// Returns an error if SMS is not configured.
func SendSMS(ctx context.Context, phone, content string) error {
	cfg := loadSMSConfig(ctx)
	if cfg == nil || cfg.Provider == "" {
		return fmt.Errorf("短信服务未配置")
	}
	p, ok := smsprovider.Get(cfg.Provider)
	if !ok {
		return fmt.Errorf("未知短信提供商: %s", cfg.Provider)
	}
	return p.Send(ctx, smsprovider.Config{
		Provider:        cfg.Provider,
		AccessKeyID:     cfg.AccessKeyID,
		AccessKeySecret: cfg.AccessKeySecret,
		SignName:        cfg.SignName,
		TemplateCode:    cfg.TemplateCode,
	}, phone, content)
}
