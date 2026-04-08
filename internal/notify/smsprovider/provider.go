package smsprovider

import "context"

// SMSProvider is the interface every SMS backend must implement.
type SMSProvider interface {
	Name() string
	Send(ctx context.Context, cfg Config, phone, content string) error
}

// Config is the shared SMS config loaded from the options table.
type Config struct {
	Provider        string
	AccessKeyID     string
	AccessKeySecret string
	SignName        string
	TemplateCode    string
}

var registry = map[string]SMSProvider{}

func Register(p SMSProvider) {
	registry[p.Name()] = p
}

func Get(name string) (SMSProvider, bool) {
	p, ok := registry[name]
	return p, ok
}
