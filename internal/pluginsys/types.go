package pluginsys

import (
	"context"
	"strings"
)

// ─── SettingType ─────────────────────────────────────────────────────────────

// SettingType enumerates valid field types for plugin settings.
type SettingType string

const (
	SettingTypeString   SettingType = "string"
	SettingTypePassword SettingType = "password"
	SettingTypeNumber   SettingType = "number"
	SettingTypeBoolean  SettingType = "boolean"
	SettingTypeSelect   SettingType = "select"
	SettingTypeTextarea SettingType = "textarea"
)

// ─── Filter event constants ───────────────────────────────────────────────────

const (
	FilterPostCreate    = "filter:post.create"
	FilterPostUpdate    = "filter:post.update"
	FilterPostDelete    = "filter:post.delete"
	FilterPostPublish   = "filter:post.publish"   // v2: fires before status changes to published
	FilterPostRestore   = "filter:post.restore"   // v2: fires before restoring from trash
	FilterCommentCreate = "filter:comment.create"
	FilterCommentDelete = "filter:comment.delete"
	FilterCommentUpdate = "filter:comment.update" // v2: fires before a comment is edited
	FilterTermCreate    = "filter:term.create"
	FilterUserRegister  = "filter:user.register"
	FilterUserUpdate    = "filter:user.update"
	FilterUserLogin     = "filter:user.login"     // v2: fires before login; abort to block
	FilterMediaUpload   = "filter:media.upload"
	FilterContentRender = "filter:content.render"
)

// ─── Manifest ────────────────────────────────────────────────────────────────

// SettingField describes one admin-configurable plugin parameter.
type SettingField struct {
	Key         string      `json:"key"`
	Label       string      `json:"label"`
	Type        SettingType `json:"type"`
	Required    bool        `json:"required"`
	Default     any         `json:"default"`
	Placeholder string      `json:"placeholder"`
	Description string      `json:"description"`
	Options     []string    `json:"options"`  // used when Type == SettingTypeSelect
	Group       string      `json:"group,omitempty"`   // v2: groups settings into collapsible sections
	ShowIf      string      `json:"showIf,omitempty"`  // v2: JS expression; field hidden when falsy, e.g. "advanced_mode === true"
}

// PluginCapabilities declares which platform APIs the plugin is permitted to use.
type PluginCapabilities struct {
	HTTP   *HTTPCap   `json:"http,omitempty"`
	Store  *StoreCap  `json:"store,omitempty"`
	Events *EventsCap `json:"events,omitempty"`
	DB     bool       `json:"db,omitempty"`     // Phase 4.1: access to plugin-prefixed tables
	AI     bool       `json:"ai,omitempty"`     // Phase 5.1: access to nuxtblog.ai service
}

// HTTPCap grants access to HTTP fetch.
type HTTPCap struct {
	// Allow lists domains the plugin may contact. Empty = any domain allowed.
	Allow []string `json:"allow"`
	// TimeoutMs overrides the default 15-second per-request timeout.
	// 0 means use the default.
	TimeoutMs int `json:"timeout_ms"`
}

// StoreCap controls store access granularity.
type StoreCap struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}

// EventsCap restricts which events the plugin may subscribe to.
// Patterns support a single trailing "*" wildcard (e.g. "post.*").
type EventsCap struct {
	Subscribe []string `json:"subscribe"`
}

// WebhookDef declares one outbound webhook registered by a plugin.
// When a matching event fires, the platform POSTs the event payload as JSON to URL.
//
// Both URL and header values support {{settings.key}} placeholders which are
// resolved at dispatch time using the plugin's settings cache (30-second TTL).
// Use this for tokens and URLs that admins configure through the settings UI —
// never hardcode secrets directly in the manifest.
//
// Example:
//
//	{ "url": "{{settings.webhook_url}}", "headers": { "Authorization": "Bearer {{settings.token}}" } }
type WebhookDef struct {
	// URL is the endpoint to POST to. Supports {{settings.key}} interpolation.
	URL string `json:"url"`
	// Events is a list of event names or patterns to subscribe to.
	// Patterns support a single trailing "*" wildcard (e.g. "post.*").
	// Use "*" to match every event.
	Events []string `json:"events"`
	// Headers are extra HTTP headers sent with every request.
	// Values support {{settings.key}} interpolation.
	Headers map[string]string `json:"headers,omitempty"`
}

// isEventMatch reports whether pattern matches eventName.
// Supported patterns: exact match, "post.*" prefix wildcard, or "*" (any).
func isEventMatch(pattern, eventName string) bool {
	if pattern == "*" || pattern == eventName {
		return true
	}
	if strings.HasSuffix(pattern, ".*") {
		prefix := strings.TrimSuffix(pattern, "*") // e.g. "post."
		return strings.HasPrefix(eventName, prefix)
	}
	return false
}

// HostVersion is the current plugin API version exposed by this build.
// Plugins may declare minHostVersion to require a minimum version.
const HostVersion = "2.0.0"

// TrustLevel controls how the frontend sandbox runs admin_js / public_js.
type TrustLevel string

const (
	TrustLevelOfficial  TrustLevel = "official"  // runs in main page context, full nuxtblogAdmin API
	TrustLevelCommunity TrustLevel = "community" // sandboxed iframe + postMessage, restricted API
	TrustLevelLocal     TrustLevel = "local"     // user-installed, runs in main context, user responsibility
)

// Manifest holds the plugin metadata parsed from package.json's "plugin" field.
type Manifest struct {
	Name         string             `json:"name"`
	Title        string             `json:"title"`
	Description  string             `json:"description"`
	Version      string             `json:"version"`
	Author       string             `json:"author"`
	Icon         string             `json:"icon"`
	Entry        string             `json:"entry"`      // path inside the archive, e.g. "dist/index.js"
	Settings     []SettingField     `json:"settings"`   // admin-configurable parameters
	CSS          string             `json:"css"`        // optional CSS injected into frontend <head>
	Priority     int                `json:"priority"`   // execution order: lower runs first (default 10)
	Capabilities PluginCapabilities `json:"capabilities"`
	Webhooks     []WebhookDef       `json:"webhooks,omitempty"`  // W-B1
	Pipelines    []PipelineDef      `json:"pipelines,omitempty"` // P-B1

	// v2 fields
	MinHostVersion   string     `json:"minHostVersion,omitempty"`   // e.g. "2.0.0"; engine rejects if host is older
	TrustLevel       TrustLevel `json:"trust_level,omitempty"`      // frontend sandbox level; default "community"
	ActivationEvents []string   `json:"activationEvents,omitempty"` // lazy-load triggers; nil/empty = "onStartup"
	AdminJS          string     `json:"admin_js,omitempty"`         // browser-side script for admin panel
	PublicJS         string     `json:"public_js,omitempty"`        // browser-side script for public frontend
	Routes           []RouteDef    `json:"routes,omitempty"`           // Phase 2.7: custom HTTP endpoints
	Contributes      *Contributes  `json:"contributes,omitempty"`    // Phase 2.2: UI contribution points
	Migrations       []MigrationDef `json:"migrations,omitempty"`   // Phase 4.1: DB schema migrations
	Pages            []PageDef     `json:"pages,omitempty"`          // Phase 4.2: frontend route extensions
	Service          *ServiceDef   `json:"service,omitempty"`        // Phase 4.3: external service proxy
	Type             string        `json:"type,omitempty"`           // "builtin" | "js" | "yaml" | "full"
}

// ─── Route definitions (Phase 2.7) ──────────────────────────────────────────

// RouteDef declares one custom HTTP endpoint registered by a plugin.
// Paths are forced to /api/plugin/{id}/ prefix to prevent route collisions.
type RouteDef struct {
	Method      string `json:"method"`                // HTTP method: GET, POST, PUT, DELETE
	Path        string `json:"path"`                  // e.g. "/api/plugin/ai-polish/invoke"
	Fn          string `json:"fn"`                    // exported JS function name to call
	Auth        string `json:"auth"`                  // "admin", "user", or "public"
	TimeoutMs   int    `json:"timeout_ms,omitempty"`  // per-request timeout; 0 = 15000ms
	Description string `json:"description,omitempty"` // OpenAPI description
}

// ─── Contribution Points (Phase 2.2) ────────────────────────────────────────

// Contributes declares UI extension points a plugin provides.
type Contributes struct {
	Commands   []CommandDef            `json:"commands,omitempty"`
	Navigation []NavigationDef         `json:"navigation,omitempty"`
	Menus      map[string][]MenuEntry  `json:"menus,omitempty"`
	Views      map[string][]ViewDef    `json:"views,omitempty"`
}

// CommandDef declares a command that can be triggered from menus or keyboard shortcuts.
type CommandDef struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	TitleEn  string `json:"title_en,omitempty"`
	Icon     string `json:"icon,omitempty"`
	Shortcut string `json:"shortcut,omitempty"`
}

// NavigationDef declares a navigation item in the admin sidebar or topbar.
type NavigationDef struct {
	Slot  string `json:"slot"`            // e.g. "admin:sidebar-nav"
	Title string `json:"title"`
	Icon  string `json:"icon,omitempty"`
	Route string `json:"route"`           // e.g. "/admin/ai"
	Order int    `json:"order,omitempty"` // sort order within the slot
}

// MenuEntry references a command by ID and optionally restricts visibility.
type MenuEntry struct {
	Command string `json:"command"` // references CommandDef.ID
	When    string `json:"when,omitempty"`
}

// ViewDef declares a panel or widget in a named slot.
type ViewDef struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type,omitempty"` // "webview" or empty for declarative
	Icon  string `json:"icon,omitempty"`
}

// ─── Migration definitions (Phase 4.1) ──────────────────────────────────────

// MigrationDef declares one versioned schema migration for a plugin.
// Table names must be prefixed with plugin_{sanitized_id}_.
type MigrationDef struct {
	Version int    `json:"version"`
	Up      string `json:"up"`              // DDL statement (CREATE, ALTER, CREATE INDEX)
	Down    string `json:"down,omitempty"`  // rollback DDL
}

// ─── Page definitions (Phase 4.2) ───────────────────────────────────────────

// PageDef declares a frontend page registered by a plugin.
type PageDef struct {
	Path      string  `json:"path"`               // e.g. "/admin/community" or "/community"
	Slot      string  `json:"slot"`               // "admin" or "public"
	Component string  `json:"component"`          // Vue component export name
	Title     string  `json:"title,omitempty"`
	Nav       *NavDef `json:"nav,omitempty"`       // optional sidebar nav entry
}

// NavDef is the navigation entry for a plugin page.
type NavDef struct {
	Group string `json:"group,omitempty"`
	Icon  string `json:"icon,omitempty"`
	Order int    `json:"order,omitempty"`
}

// ─── Service proxy (Phase 4.3) ──────────────────────────────────────────────

// ServiceDef declares an external service that the blog proxies requests to.
type ServiceDef struct {
	Proxy  string `json:"proxy"`  // path prefix, e.g. "/api/plugin/community/"
	Target string `json:"target"` // upstream URL, e.g. "http://localhost:8090"
}

// ─── PluginCtx ───────────────────────────────────────────────────────────────

// PluginCtx is the context object passed to every filter handler.
type PluginCtx struct {
	// Context carries the request deadline and trace information.
	Context context.Context
	// Event is the filter event name (e.g. "filter:post.create").
	Event string
	// Input is a deep-copy snapshot of Data taken before the chain starts.
	// It is read-only and used for diff/audit logging.
	Input map[string]any
	// Data is the mutable payload. Plugins read and write this field.
	Data map[string]any
	// Meta is a request-scoped KV store for inter-plugin communication.
	// Values written by earlier plugins are visible to later ones.
	Meta    map[string]any
	aborted bool
	reason  string
}

// Next explicitly signals that this handler is done and the chain should continue.
// Calling Next is optional — the chain always continues unless Abort is called.
func (c *PluginCtx) Next() {}

// Abort stops the filter chain immediately. All subsequent plugin handlers are
// skipped. The caller of Filter() receives an error wrapping reason.
func (c *PluginCtx) Abort(reason string) {
	c.aborted = true
	c.reason = reason
}

// IsAborted reports whether Abort was called by any handler in the chain.
func (c *PluginCtx) IsAborted() bool { return c.aborted }

// ─── Pipeline ────────────────────────────────────────────────────────────────

// PipelineDef declares one multi-step async workflow in a plugin manifest.
type PipelineDef struct {
	Name    string    `json:"name"`    // unique within the plugin, used in logs
	Trigger string    `json:"trigger"` // event name or pattern (e.g. "post.*")
	Steps   []StepDef `json:"steps"`
}

// StepDef describes one step within a pipeline.
// Type selects which executor runs it: "js", "webhook", or "condition".
type StepDef struct {
	Type string `json:"type"` // "js" | "webhook" | "condition"
	Name string `json:"name"` // human-readable label for logs and UI

	// type = js: call an exported JS function by name
	Fn string `json:"fn,omitempty"`

	// type = webhook: POST the StepContext.Data as JSON
	URL     string            `json:"url,omitempty"`
	Headers map[string]string `json:"headers,omitempty"`

	// type = condition: evaluate a JS boolean expression, then run Then or Else
	If   string    `json:"if,omitempty"`
	Then []StepDef `json:"then,omitempty"`
	Else []StepDef `json:"else,omitempty"`

	// Reliability
	TimeoutMs int `json:"timeout_ms,omitempty"` // per-step timeout; 0 = 5000 ms
	Retry     int `json:"retry,omitempty"`       // extra attempts on failure; 0 = none
}

// StepContext flows through every step of a pipeline.
// Unlike PluginCtx it allows HTTP and has no Input snapshot.
type StepContext struct {
	Context   context.Context
	EventName string
	Data      map[string]any // shared mutable payload
	Meta      map[string]any // step-to-step KV store
	aborted   bool
	reason    string
}

// Abort stops the pipeline immediately. Subsequent steps are skipped.
func (c *StepContext) Abort(reason string) { c.aborted = true; c.reason = reason }

// IsAborted reports whether Abort was called.
func (c *StepContext) IsAborted() bool { return c.aborted }
