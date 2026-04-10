package v1

import "github.com/gogf/gf/v2/frame/g"

// AIConfig is a single AI provider configuration stored in the options table.
type AIConfig struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	APIFormat string `json:"api_format"`        // "openai" (default) | "claude"
	Label     string `json:"label,omitempty"`   // optional display label, e.g. "OpenAI", "My Proxy"
	APIKey    string `json:"api_key"`           // masked on read
	Model     string `json:"model"`
	BaseURL   string `json:"base_url"`          // required for custom endpoints
	TimeoutMs int    `json:"timeout_ms"`
	IsActive  bool   `json:"is_active"`
}

// ── Config CRUD ───────────────────────────────────────────────────────────────

type AIListConfigsReq struct {
	g.Meta `path:"/admin/ai/configs" method:"get" tags:"AI" summary:"List AI configs" auth:"true"`
}
type AIListConfigsRes struct {
	Items    []AIConfig `json:"items"`
	ActiveID string     `json:"active_id"`
}

type AICreateConfigReq struct {
	g.Meta    `path:"/admin/ai/configs" method:"post" tags:"AI" summary:"Create AI config" auth:"true"`
	Name      string `json:"name"       v:"required" dc:"Display name, e.g. My GPT-4o"`
	APIFormat string `json:"api_format" v:"in:openai,claude" dc:"openai (default, OpenAI-compatible) or claude (Anthropic format)"`
	Label     string `json:"label"      dc:"Optional: provider label shown in UI"`
	APIKey    string `json:"api_key"    dc:"API key; leave empty for local services"`
	Model     string `json:"model"      v:"required" dc:"Model name, e.g. gpt-4o-mini"`
	BaseURL   string `json:"base_url"   dc:"API base URL, e.g. https://ai.071129.xyz"`
	TimeoutMs int    `json:"timeout_ms" dc:"Timeout ms, default 30000"`
}
type AICreateConfigRes struct {
	Item AIConfig `json:"item"`
}

type AIUpdateConfigReq struct {
	g.Meta    `path:"/admin/ai/configs/{id}" method:"put" tags:"AI" summary:"Update AI config" auth:"true"`
	ID        string `json:"-"          v:"required" dc:"config id"`
	Name      string `json:"name"`
	APIFormat string `json:"api_format" v:"in:openai,claude"`
	Label     string `json:"label"`
	APIKey    string `json:"api_key"    dc:"New key; send masked value to keep existing"`
	Model     string `json:"model"`
	BaseURL   string `json:"base_url"`
	TimeoutMs int    `json:"timeout_ms"`
}
type AIUpdateConfigRes struct {
	Item AIConfig `json:"item"`
}

type AIDeleteConfigReq struct {
	g.Meta `path:"/admin/ai/configs/{id}" method:"delete" tags:"AI" summary:"Delete AI config" auth:"true"`
	ID     string `json:"-" v:"required" dc:"config id"`
}
type AIDeleteConfigRes struct{}

type AIActivateConfigReq struct {
	g.Meta `path:"/admin/ai/configs/{id}/activate" method:"post" tags:"AI" summary:"Set active AI config" auth:"true"`
	ID     string `json:"-" v:"required" dc:"config id"`
}
type AIActivateConfigRes struct{}

type AITestConfigReq struct {
	g.Meta `path:"/admin/ai/test" method:"post" tags:"AI" summary:"Test AI config" auth:"true"`
	ID     string `json:"id" v:"required" dc:"config id to test"`
}
type AITestConfigRes struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// ── AI Actions ────────────────────────────────────────────────────────────────

type AIFromURLReq struct {
	g.Meta `path:"/ai/from-url" method:"post" tags:"AI" summary:"Generate article from URL" auth:"true"`
	URL    string `json:"url"   v:"required|url" dc:"Web page URL"`
	Style  string `json:"style" dc:"casual|formal|concise (optional)"`
}
type AIFromURLRes struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}

