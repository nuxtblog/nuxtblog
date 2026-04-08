package consts

// Notify channel type constants — a channel's Name() must return one of these.
// To add a new channel: implement notify.Channel, register via notify.Register() in init(),
// then add an entry below.
const (
	NotifyChannelInApp   = "in_app"   // built-in persistent in-app notification
	NotifyChannelEmail   = "email"    // SMTP email
	NotifyChannelSMS     = "sms"      // SMS via provider (Aliyun / Tencent / Twilio …)
	NotifyChannelWebhook = "webhook"  // generic outbound HTTP webhook
	NotifyChannelWechat  = "wechat"   // WeChat official account template message
)

// NotifyChannelDef describes a notification delivery channel.
type NotifyChannelDef struct {
	Slug      string // must match Channel.Name()
	LabelZh   string
	LabelEn   string
	ConfigKey string // options DB key for admin-editable settings; "" = no DB config needed
}

// BuiltinNotifyChannels lists every channel implemented in this codebase.
// The admin settings page reads this to show which channels can be configured.
// Adding a channel requires: (1) implement notify.Channel + register via init(),
// (2) add a ConfigKey entry in migrate.go defaultOpts if needed,
// (3) add an entry here.
var BuiltinNotifyChannels = []NotifyChannelDef{
	{Slug: NotifyChannelInApp,   LabelZh: "站内通知", LabelEn: "In-App",       ConfigKey: ""},
	{Slug: NotifyChannelEmail,   LabelZh: "邮件通知", LabelEn: "Email",        ConfigKey: "notify_email"},
	{Slug: NotifyChannelSMS,     LabelZh: "短信通知", LabelEn: "SMS",          ConfigKey: "notify_sms"},
	{Slug: NotifyChannelWebhook, LabelZh: "Webhook",  LabelEn: "Webhook",      ConfigKey: "notify_webhook"},
	{Slug: NotifyChannelWechat,  LabelZh: "微信通知", LabelEn: "WeChat",       ConfigKey: "notify_wechat"},
}
