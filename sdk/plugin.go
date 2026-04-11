// Package plugin defines the Go-native plugin SDK for nuxtblog.
//
// Plugins implement the Plugin interface and optionally implement
// additional interfaces (HasRoutes, HasEvents, etc.) for extra capabilities.
//
// Built-in plugins register via init() + Register().
// Third-party plugins are loaded as JavaScript via Goja.
package sdk

import (
	"context"
	"net/http"
)

// Plugin is the interface every Go plugin must implement.
//
// It includes all lifecycle methods with no-op defaults available via BasePlugin.
// Embed BasePlugin in your struct and override only the methods you need.
//
// The wide interface ensures that the plugin manager can call all methods
// uniformly without type-asserting optional capabilities.
type Plugin interface {
	Manifest() Manifest
	Activate(ctx PluginContext) error
	Deactivate() error
	Filters() []FilterDef
	Routes(r RouteRegistrar)
	OnEvent(ctx context.Context, event string, data map[string]any)
	Migrations() []Migration
}

// BasePlugin provides no-op defaults for all optional Plugin methods.
// Embed this in your plugin struct and override only the methods you need:
//
//	type MyPlugin struct { sdk.BasePlugin }
//	func (p *MyPlugin) Manifest() sdk.Manifest { ... }
//	func (p *MyPlugin) Filters() []sdk.FilterDef { ... } // only override what you need
type BasePlugin struct{}

func (BasePlugin) Activate(PluginContext) error                          { return nil }
func (BasePlugin) Deactivate() error                                     { return nil }
func (BasePlugin) Filters() []FilterDef                                  { return nil }
func (BasePlugin) Routes(RouteRegistrar)                                 {}
func (BasePlugin) OnEvent(context.Context, string, map[string]any)       {}
func (BasePlugin) Migrations() []Migration                               { return nil }

// ─── Optional capability marker interfaces (used by builtin plugin manager) ──

// HasRoutes indicates the plugin registers custom HTTP endpoints.
type HasRoutes interface {
	Routes(r RouteRegistrar)
}

// HasEvents indicates the plugin subscribes to platform events.
type HasEvents interface {
	OnEvent(ctx context.Context, event string, data map[string]any)
}

// HasFilters indicates the plugin intercepts content mutations.
type HasFilters interface {
	Filters() []FilterDef
}

// HasMigrations indicates the plugin needs database schema setup.
type HasMigrations interface {
	Migrations() []Migration
}

// HasActivate indicates the plugin needs initialization with platform services.
type HasActivate interface {
	Activate(ctx PluginContext) error
}

// HasDeactivate indicates the plugin needs cleanup on shutdown.
type HasDeactivate interface {
	Deactivate() error
}

// HasAssets indicates the plugin provides embedded frontend assets
// (admin.mjs, public.mjs, etc.) that should be written to data/plugins/{id}/.
type HasAssets interface {
	Assets() map[string][]byte // filename -> content
}

// ─── Plugin Types ─────��───────────────────────────────────────────────────

// PluginType enumerates plugin runtime types.
const (
	TypeBuiltin = "builtin" // compiled into binary, installable from marketplace (requires server restart)
	TypeJS      = "js"      // JavaScript source interpreted by Goja, installable
	TypeYAML    = "yaml"    // declarative YAML (webhooks, filters), installable
	TypeUI      = "ui"      // frontend-only (admin.mjs/public.mjs), installable
	TypeFull    = "full"    // JS + frontend assets, installable
)

// ─── Manifest ��────────────────────────────────────────────────────────────

// Manifest is the unified plugin configuration.
// ALL plugin types (builtin, go, yaml, ui, full) share this single struct.
// It maps 1:1 to plugin.yaml on disk.
type Manifest struct {
	// ── Core identity (required) ──
	ID          string `yaml:"id"          json:"id"`
	Title       string `yaml:"title"       json:"title"`
	Version     string `yaml:"version"     json:"version"`
	Icon        string `yaml:"icon"        json:"icon"`
	Author      string `yaml:"author"      json:"author"`
	Description string `yaml:"description" json:"description"`

	// ── Type & trust ──
	Type       string `yaml:"type"        json:"type"`                   // builtin, go, yaml, ui, full
	TrustLevel string `yaml:"trust_level" json:"trust_level,omitempty"`
	License    string `yaml:"license"     json:"license,omitempty"`
	SDKVersion string `yaml:"sdk_version" json:"sdk_version,omitempty"`
	Priority   int    `yaml:"priority"    json:"priority,omitempty"`     // execution order; 0 = default 10

	// ── Runtime & distribution ──
	Runtime string `yaml:"runtime" json:"runtime,omitempty"` // compiled | interpreted
	Bundled bool   `yaml:"bundled" json:"bundled,omitempty"` // included in official prebuilt binary

	// ── Marketplace (for installable types) ──
	Repo     string   `yaml:"repo"     json:"repo,omitempty"`
	Homepage string   `yaml:"homepage" json:"homepage,omitempty"`
	Tags     []string `yaml:"tags"     json:"tags,omitempty"`

	// ── JS plugin specific ──
	JSEntry string `yaml:"js_entry" json:"js_entry,omitempty"` // e.g. "plugin.js"; default "plugin.js"

	// ── Frontend assets ──
	AdminJS  string `yaml:"admin_js"  json:"admin_js,omitempty"`
	PublicJS string `yaml:"public_js" json:"public_js,omitempty"`
	CSS      string `yaml:"css"       json:"css,omitempty"`

	// ── Configuration schema ──
	Settings []SettingDef `yaml:"settings" json:"settings,omitempty"`

	// ── Dependencies ──
	Depends []Dependency `yaml:"depends" json:"depends,omitempty"`

	// ── Functionality declarations ──
	Pages       []PageDef      `yaml:"pages"       json:"pages,omitempty"`
	Routes      []RouteDef     `yaml:"routes"      json:"routes,omitempty"`
	Migrations  []Migration    `yaml:"migrations"  json:"migrations,omitempty"`
	Contributes *Contributes   `yaml:"contributes" json:"contributes,omitempty"`

	// ── Capabilities ──
	Capabilities *Capabilities `yaml:"capabilities" json:"capabilities,omitempty"`

	// ── Permissions (frontend API access) ──
	Permissions []string `yaml:"permissions" json:"permissions,omitempty"`

	// ── YAML plugin specific (declarative logic) ──
	Webhooks []WebhookDef `yaml:"webhooks" json:"webhooks,omitempty"`
	Filters  []YAMLFilter `yaml:"filters"  json:"filters,omitempty"`
}

// ─── Capabilities ────────────────────────────────────────────────────────

// DBCapability declares the database access level for a plugin.
type DBCapability struct {
	Own    bool            `yaml:"own"    json:"own,omitempty"`
	Tables []DBTableAccess `yaml:"tables" json:"tables,omitempty"`
	Raw    bool            `yaml:"raw"    json:"raw,omitempty"`
}

// DBTableAccess declares permitted operations on a specific core table.
type DBTableAccess struct {
	Table string   `yaml:"table" json:"table"`
	Ops   []string `yaml:"ops"   json:"ops"`
}

// Capabilities declares which platform APIs the plugin is permitted to use.
type Capabilities struct {
	DB *DBCapability `yaml:"db" json:"db,omitempty"`
}

// ─── Setting ──────────────────────────────────────────────────────────────

// SettingDef describes an admin-configurable setting field.
type SettingDef struct {
	Key         string   `yaml:"key"         json:"key"`
	Label       string   `yaml:"label"       json:"label"`
	Type        string   `yaml:"type"        json:"type"` // "string", "number", "boolean", "select", "password", "textarea"
	Required    bool     `yaml:"required"    json:"required,omitempty"`
	Default     any      `yaml:"default"     json:"default,omitempty"`
	Placeholder string   `yaml:"placeholder" json:"placeholder,omitempty"`
	Description string   `yaml:"description" json:"description,omitempty"`
	Options     []string `yaml:"options"     json:"options,omitempty"`
	Group       string   `yaml:"group"       json:"group,omitempty"`
	Shared      bool     `yaml:"shared"      json:"shared,omitempty"`
}

// ─── Page ─────────���───────────────────────────────────────────────────────

// PageDef declares a frontend page registered by the plugin.
type PageDef struct {
	Path      string  `yaml:"path"      json:"path"`
	Slot      string  `yaml:"slot"      json:"slot"`      // "admin" or "public"
	Component string  `yaml:"component" json:"component"` // Vue component export name
	Title     string  `yaml:"title"     json:"title,omitempty"`
	Nav       *NavDef `yaml:"nav"       json:"nav,omitempty"`
}

// NavDef is the navigation entry for a plugin page.
type NavDef struct {
	Group string `yaml:"group" json:"group,omitempty"`
	Icon  string `yaml:"icon"  json:"icon,omitempty"`
	Order int    `yaml:"order" json:"order,omitempty"`
}

// ─── Route ────────────��───────────────────────────────────────────────────

// RouteDef declares one custom HTTP endpoint in plugin.yaml.
type RouteDef struct {
	Method      string `yaml:"method"      json:"method"`
	Path        string `yaml:"path"        json:"path"`
	Auth        string `yaml:"auth"        json:"auth,omitempty"`        // "admin", "user", "public"
	Fn          string `yaml:"fn"          json:"fn,omitempty"`          // exported function name (Goja)
	Description string `yaml:"description" json:"description,omitempty"`
}

// ─── Contributes ────────────��─────────────────────────────────────────────

// Contributes declares UI extension points a plugin provides.
type Contributes struct {
	Commands   []CommandDef           `yaml:"commands"   json:"commands,omitempty"`
	Menus      map[string][]MenuEntry `yaml:"menus"      json:"menus,omitempty"`
	Navigation []NavigationDef        `yaml:"navigation" json:"navigation,omitempty"`
	Views      map[string][]ViewDef   `yaml:"views"      json:"views,omitempty"`
}

// NavigationDef declares a navigation item injected into a named slot.
type NavigationDef struct {
	Slot   string `yaml:"slot"   json:"slot"`
	Title  string `yaml:"title"  json:"title"`
	Icon   string `yaml:"icon"   json:"icon,omitempty"`
	Route  string `yaml:"route"  json:"route"`
	Order  int    `yaml:"order"  json:"order,omitempty"`
	Parent string `yaml:"parent" json:"parent,omitempty"`
}

// ViewDef declares a view panel injected into a named slot.
type ViewDef struct {
	ID    string `yaml:"id"    json:"id"`
	Title string `yaml:"title" json:"title"`
	Type  string `yaml:"type"  json:"type,omitempty"`
	Icon  string `yaml:"icon"  json:"icon,omitempty"`
}

// CommandDef declares a command triggered from menus or keyboard shortcuts.
type CommandDef struct {
	ID      string `yaml:"id"       json:"id"`
	Title   string `yaml:"title"    json:"title"`
	TitleEn string `yaml:"title_en" json:"title_en,omitempty"`
	Icon    string `yaml:"icon"     json:"icon,omitempty"`
}

// MenuEntry references a command by ID.
type MenuEntry struct {
	Command string `yaml:"command" json:"command"`
}

// ─── Webhook ──────────────────────────────────────────────────────────────

// WebhookDef declares one outbound webhook. URL and Headers support
// {{settings.key}} interpolation resolved at dispatch time.
type WebhookDef struct {
	URL     string            `yaml:"url"     json:"url"`
	Events  []string          `yaml:"events"  json:"events"`
	Headers map[string]string `yaml:"headers" json:"headers,omitempty"`
}

// ─── YAML Filter ──────────────────────────────────────────────────────────

// YAMLFilter is a declarative filter with simple matching rules.
type YAMLFilter struct {
	Event string     `yaml:"event" json:"event"`
	Rules []YAMLRule `yaml:"rules" json:"rules"`
}

// YAMLRule is a single validation rule applied to a data field.
type YAMLRule struct {
	Field        string   `yaml:"field"                    json:"field"`
	MinLength    int      `yaml:"min_length,omitempty"     json:"min_length,omitempty"`
	MaxLength    int      `yaml:"max_length,omitempty"     json:"max_length,omitempty"`
	BlockedWords []string `yaml:"blocked_words,omitempty"  json:"blocked_words,omitempty"`
	Regex        string   `yaml:"regex,omitempty"          json:"regex,omitempty"`
	NotRegex     string   `yaml:"not_regex,omitempty"      json:"not_regex,omitempty"`
	Message      string   `yaml:"message"                  json:"message"`
}

// RouteRegistrar is used by plugins to declare HTTP routes.
type RouteRegistrar interface {
	Handle(method, path string, handler http.HandlerFunc, opts ...RouteOption)
}

// RouteOption configures a route.
type RouteOption func(*RouteConfig)

// RouteConfig holds route configuration set by RouteOption functions.
type RouteConfig struct {
	Auth string // "admin", "user", "public"
}

// ApplyOptions creates a RouteConfig from a list of options.
func ApplyOptions(opts []RouteOption) *RouteConfig {
	cfg := &RouteConfig{}
	for _, o := range opts {
		o(cfg)
	}
	return cfg
}

// WithAuth sets the authentication requirement for a route.
func WithAuth(auth string) RouteOption {
	return func(c *RouteConfig) {
		c.Auth = auth
	}
}

// FilterDef describes a content filter.
type FilterDef struct {
	Event   string                    // e.g. "content.render", "post.create"
	Handler func(ctx *FilterContext)
}

// FilterContext is passed to filter handlers.
type FilterContext struct {
	Context context.Context
	Event   string
	Data    map[string]any
	Meta    map[string]any
	aborted bool
	reason  string
}

// Abort stops the filter chain.
func (c *FilterContext) Abort(reason string) {
	c.aborted = true
	c.reason = reason
}

// IsAborted reports whether Abort was called.
func (c *FilterContext) IsAborted() bool { return c.aborted }

// AbortReason returns the abort reason.
func (c *FilterContext) AbortReason() string { return c.reason }

// Migration describes a versioned schema migration.
type Migration struct {
	Version int    `yaml:"version" json:"version"`
	Up      string `yaml:"up"      json:"up"`
	Down    string `yaml:"down"    json:"down,omitempty"`
}

// Dependency declares a requirement on another plugin.
type Dependency struct {
	ID       string `yaml:"id"       json:"id"`
	Version  string `yaml:"version"  json:"version,omitempty"`  // semver constraint, e.g. ">=1.0.0"
	Optional bool   `yaml:"optional" json:"optional,omitempty"`
}

// PluginQuery provides runtime queries about other plugins.
type PluginQuery interface {
	IsAvailable(id string) bool
	GetVersion(id string) string
	GetSetting(pluginID, key string) (any, error)
}

// PluginContext provides platform services to plugins during activation.
type PluginContext struct {
	DB       DB
	Store    Store
	Settings Settings
	Log      Logger
	Plugins  PluginQuery
	AI       AI
}

// DB provides isolated database access for plugins.
type DB interface {
	Query(sql string, args ...any) ([]map[string]any, error)
	Execute(sql string, args ...any) (int64, error)
}

// Store provides key-value storage for plugins.
type Store interface {
	Get(key string) (any, error)
	Set(key string, value any) error
	Delete(key string) error
	Increment(key string, delta ...int64) (int64, error)
	DeletePrefix(prefix string) (int64, error)
}

// Settings provides read-only access to plugin configuration.
type Settings interface {
	Get(key string) any
	GetAll() map[string]any
}

// Logger provides structured logging for plugins.
type Logger interface {
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	Debug(msg string)
}

// ── AI types ────────────────────────────────────────────────────────────────

// Role identifies the author of a chat message.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// Message is a single message in a conversation.
type Message struct {
	Role    Role
	Content string
}

// Messages is a convenience constructor: Messages(RoleSystem, "sys", RoleUser, "hi").
func Messages(args ...any) []Message {
	var msgs []Message
	for i := 0; i+1 < len(args); i += 2 {
		r, _ := args[i].(Role)
		c, _ := args[i+1].(string)
		msgs = append(msgs, Message{Role: r, Content: c})
	}
	return msgs
}

// AIRequest describes a one-shot or multi-turn LLM call.
type AIRequest struct {
	Messages    []Message // conversation messages (required)
	MaxTokens   int       // 0 = provider default
	Temperature float64   // 0 = provider default
}

// AIUsage tracks token consumption.
type AIUsage struct {
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// AIResponse holds the result of an AI generation call.
type AIResponse struct {
	Text         string   // primary text output
	FinishReason string   // "stop", "length", etc.; empty if unknown
	Usage        *AIUsage // nil if provider doesn't report usage
}

// AI is the safe AI surface exposed to plugins.
// Credentials are NEVER passed through this interface.
type AI interface {
	Generate(ctx context.Context, req AIRequest) (*AIResponse, error)
}
