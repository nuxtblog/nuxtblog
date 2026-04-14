package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// SettingField mirrors engine.SettingField for API responses.
type SettingField struct {
	Key         string   `json:"key"`
	Label       string   `json:"label"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Default     any      `json:"default"`
	Placeholder string   `json:"placeholder"`
	Description string   `json:"description"`
	Options     []string `json:"options"`
	Shared      bool     `json:"shared,omitempty"`
}

// PluginItem is returned in list responses.
type PluginItem struct {
	Id           string `json:"id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	Version      string `json:"version"`
	Author       string `json:"author"`
	Icon         string `json:"icon"`
	RepoUrl      string `json:"repo_url"`
	Enabled      bool   `json:"enabled"`
	InstalledAt  string `json:"installed_at"`
	Capabilities string `json:"capabilities"`          // raw JSON, e.g. {} means full-access
	Source       string `json:"source"`                 // "builtin" (Go native, cannot uninstall) | "external" (installed via zip/github)
	Type         string `json:"type"`                   // "builtin" | "js" | "yaml" | "full"
	NeedRestart  bool   `json:"need_restart,omitempty"` // true when plugin requires server restart to activate
}

// ── List ──────────────────────────────────────────────────────────────────

type PluginListReq struct {
	g.Meta `path:"/admin/plugins" method:"get" tags:"Plugin" summary:"Admin: list plugins" auth:"true"`
}
type PluginListRes struct {
	Items []PluginItem `json:"items"`
}

// ── Install ───────────────────────────────────────────────────────────────

// PluginInstallReq installs a plugin from a GitHub repository URL.
// Accepted formats:
//
//	https://github.com/owner/repo
//	github.com/owner/repo
//	owner/repo
type PluginInstallReq struct {
	g.Meta          `path:"/admin/plugins" method:"post" tags:"Plugin" summary:"Admin: install plugin" auth:"true"`
	RepoUrl         string `v:"required" json:"repo_url" dc:"GitHub repo URL or owner/repo"`
	ExpectedVersion string `json:"expected_version,omitempty" dc:"Registry version to validate against GitHub release tag; empty disables validation"`
}
type PluginInstallRes struct {
	Item PluginItem `json:"item"`
}

// ── Uninstall ─────────────────────────────────────────────────────────────

type PluginUninstallReq struct {
	g.Meta `path:"/admin/plugins/{id}" method:"delete" tags:"Plugin" summary:"Admin: uninstall plugin" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}
type PluginUninstallRes struct {
	NeedRestart bool `json:"need_restart,omitempty"`
}

// ── Batch Uninstall ──────────────────────────────────────────────────────

type PluginBatchUninstallReq struct {
	g.Meta `path:"/admin/plugins/batch-uninstall" method:"post" tags:"Plugin" summary:"Admin: batch uninstall plugins" auth:"true"`
	Ids    []string `v:"required" json:"ids" dc:"plugin ids to uninstall"`
}
type PluginBatchUninstallRes struct {
	Succeeded   []string `json:"succeeded"`
	Failed      []string `json:"failed"`
	NeedRestart bool     `json:"need_restart,omitempty"`
}

// ── Upload ZIP ────────────────────────────────────────────────────────────

// PluginUploadZipReq installs a plugin from a local archive file.
// Supported formats: zip (all compression methods), tar.gz, tar.bz2, tar.xz, 7z, rar, …
// The archive must contain package.json (with a "plugin" field) and the entry script (default index.js).
type PluginUploadZipReq struct {
	g.Meta `path:"/admin/plugins/upload" method:"post" tags:"Plugin" summary:"Admin: install plugin from archive" mime:"multipart/form-data" auth:"true"`
	File   *ghttp.UploadFile `v:"required" dc:"Archive file (zip/tar.gz/7z/rar/…) containing package.json (with \"plugin\" field) and index.js"`
}
type PluginUploadZipRes struct {
	Item PluginItem `json:"item"`
}

// ── Settings ──────────────────────────────────────────────────────────────

type PluginGetSettingsReq struct {
	g.Meta `path:"/admin/plugins/{id}/settings" method:"get" tags:"Plugin" summary:"Admin: get plugin settings schema and values" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}
type PluginGetSettingsRes struct {
	Schema []SettingField         `json:"schema"`
	Values map[string]interface{} `json:"values"`
}

type PluginUpdateSettingsReq struct {
	g.Meta `path:"/admin/plugins/{id}/settings" method:"put" tags:"Plugin" summary:"Admin: update plugin settings values" auth:"true"`
	Id     string                 `v:"required" dc:"plugin id"`
	Values map[string]interface{} `json:"values"`
}
type PluginUpdateSettingsRes struct{}

// ── Styles (public) ───────────────────────────────────────────────────────

// PluginStylesReq returns concatenated CSS from all enabled plugins.
// No auth required — called by the frontend on every page load.
type PluginStylesReq struct {
	g.Meta `path:"/plugins/styles" method:"get" tags:"Plugin" summary:"Get CSS from all enabled plugins"`
}
type PluginStylesRes struct {
	CSS string `json:"css"`
}

// ── Client Plugins (public) ───────────────────────────────────────────────

// PluginClientListReq returns enabled plugins' client-side info (contributes,
// admin_js, trust_level) for the admin panel to render UI extensions.
// Requires admin auth so only logged-in admins load plugin scripts.
type PluginClientListReq struct {
	g.Meta `path:"/admin/plugins/client" method:"get" tags:"Plugin" summary:"Admin: list enabled plugins with client-side info" auth:"true"`
}

type PluginClientItem struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Icon        string `json:"icon"`
	Version     string `json:"version"`
	TrustLevel  string `json:"trust_level"`
	Contributes string `json:"contributes,omitempty"` // raw JSON
	Permissions string `json:"permissions,omitempty"` // raw JSON — frontend API permissions
}

type PluginClientListRes struct {
	Items []PluginClientItem `json:"items"`
}

// ── Public Client Plugins ─────────────────────────────────────────────────

// PluginPublicClientReq returns enabled plugins with public-facing contributes
// for the public frontend. No auth required.
type PluginPublicClientReq struct {
	g.Meta `path:"/plugins/client" method:"get" tags:"Plugin" summary:"Public: enabled plugins with public contributes"`
}
type PluginPublicClientRes struct {
	Items []PluginClientItem `json:"items"`
}

// ── Update (pull latest from GitHub) ──────────────────────────────────────

type PluginUpdateReq struct {
	g.Meta `path:"/admin/plugins/{id}/update" method:"post" tags:"Plugin" summary:"Admin: update plugin to latest version" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}
type PluginUpdateRes struct {
	Item PluginItem `json:"item"`
}

// ── Batch Update ─────────────────────────────────────────────────────────
type PluginBatchUpdateReq struct {
	g.Meta `path:"/admin/plugins/batch-update" method:"post" tags:"Plugin" summary:"Admin: batch update plugins" auth:"true"`
	Ids    []string `v:"required" json:"ids" dc:"plugin ids to update"`
}
type PluginBatchUpdateRes struct {
	Succeeded []string `json:"succeeded"`
	Failed    []string `json:"failed"`
}

// ── Toggle (enable / disable) ─────────────────────────────────────────────

type PluginToggleReq struct {
	g.Meta  `path:"/admin/plugins/{id}" method:"patch" tags:"Plugin" summary:"Admin: enable or disable plugin" auth:"true"`
	Id      string `v:"required" dc:"plugin id"`
	Enabled bool   `json:"enabled" dc:"true=enable  false=disable"`
}
type PluginToggleRes struct{}

// ── Marketplace ───────────────────────────────────────────────────────────

// MarketplaceItem is a plugin entry in the public registry.
type MarketplaceItem struct {
	Name        string   `json:"name"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Icon        string   `json:"icon"`
	Repo        string   `json:"repo"`
	Homepage    string   `json:"homepage"`
	Tags        []string `json:"tags"`
	Type        string   `json:"type"`
	Runtime     string   `json:"runtime"` // "compiled" (needs Go rebuild+restart) | "interpreted" (hot-loaded)
	IsOfficial  bool     `json:"is_official"`
	License     string   `json:"license"`
	PublishedAt string   `json:"published_at"`
	UpdatedAt   string   `json:"updated_at"`
}

type MarketplaceListReq struct {
	g.Meta  `path:"/admin/plugins/marketplace" method:"get" tags:"Plugin" summary:"Admin: list marketplace plugins" auth:"true"`
	Keyword string `dc:"search keyword"`
	Type    string `dc:"filter by type"`
}
type MarketplaceListRes struct {
	Items    []MarketplaceItem `json:"items"`
	SyncedAt string            `json:"synced_at"`
}

type MarketplaceSyncReq struct {
	g.Meta `path:"/admin/plugins/marketplace/sync" method:"post" tags:"Plugin" summary:"Admin: sync marketplace registry" auth:"true"`
}
type MarketplaceSyncRes struct {
	Count    int    `json:"count"`
	SyncedAt string `json:"synced_at"`
}

// ── Stats (4.5-B1 / 4.5-B3) ──────────────────────────────────────────────────

type PluginGetStatsReq struct {
	g.Meta `path:"/admin/plugins/{id}/stats" method:"get" tags:"Plugin" summary:"Admin: get plugin execution stats and sliding-window history" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}

// WindowBucket is one minute's aggregated counters in the 60-min sliding window.
type WindowBucket struct {
	Minute      string `json:"minute"` // RFC3339, truncated to the minute
	Invocations int64  `json:"invocations"`
	Errors      int64  `json:"errors"`
}

type PluginGetStatsRes struct {
	PluginID      string         `json:"plugin_id"`
	Invocations   int64          `json:"invocations"`
	Errors        int64          `json:"errors"`
	AvgDurationMs float64        `json:"avg_duration_ms"`
	LastError     string         `json:"last_error,omitempty"`
	LastErrorAt   string         `json:"last_error_at,omitempty"` // RFC3339, omitted when zero
	History       []WindowBucket `json:"history"`
}

// ── Errors (4.5-B2) ───────────────────────────────────────────────────────────

type PluginGetErrorsReq struct {
	g.Meta `path:"/admin/plugins/{id}/errors" method:"get" tags:"Plugin" summary:"Admin: get plugin recent error log" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}

type PluginErrorEntry struct {
	At        string `json:"at"`             // RFC3339
	EventName string `json:"event"`
	Message   string `json:"message"`
	InputDiff string `json:"input_diff,omitempty"`
}

type PluginGetErrorsRes struct {
	Items []PluginErrorEntry `json:"items"`
}

// ── Unload Impact ─────────────────────────────────────────────────────────────

type PluginUnloadImpactReq struct {
	g.Meta `path:"/admin/plugins/{id}/unload-impact" method:"get" tags:"Plugin" summary:"Admin: preview cascade unload impact" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}
type PluginUnloadImpactRes struct {
	WillUnload    []string `json:"will_unload"`
	HasDB         bool     `json:"has_db"`          // has migrations (uninstall will DROP tables)
	HasMediaCats  bool     `json:"has_media_cats"`  // has media_categories
	MediaCatSlugs []string `json:"media_cat_slugs"` // specific category slug list
}

// ── Manifest (P-B11) ──────────────────────────────────────────────────────────

type PluginGetManifestReq struct {
	g.Meta `path:"/admin/plugins/{id}/manifest" method:"get" tags:"Plugin" summary:"Admin: get plugin manifest JSON" auth:"true"`
	Id     string `v:"required" dc:"plugin id"`
}
type PluginGetManifestRes struct {
	Manifest string `json:"manifest"` // raw JSON of the full plugin manifest
}

type PluginUpdateManifestReq struct {
	g.Meta   `path:"/admin/plugins/{id}/manifest" method:"put" tags:"Plugin" summary:"Admin: update plugin manifest JSON" auth:"true"`
	Id       string `v:"required" dc:"plugin id"`
	Manifest string `json:"manifest" v:"required" dc:"raw JSON manifest to store"`
}
type PluginUpdateManifestRes struct{}

// ── Preview ───────────────────────────────────────────────────────────────

// PluginPreviewReq fetches the plugin manifest from GitHub without installing.
type PluginPreviewReq struct {
	g.Meta `path:"/admin/plugins/preview" method:"get" tags:"Plugin" summary:"Admin: preview plugin manifest from GitHub" auth:"true"`
	Repo   string `v:"required" dc:"GitHub repo URL or owner/repo"`
}

type HTTPCapPreview struct {
	Allow     []string `json:"allow"`
	TimeoutMs int      `json:"timeout_ms"`
}

type StoreCapPreview struct {
	Read  bool `json:"read"`
	Write bool `json:"write"`
}

type EventsCapPreview struct {
	Subscribe []string `json:"subscribe"`
}

type DBCapPreview struct {
	Own    bool              `json:"own,omitempty"`
	Tables []DBTablePreview  `json:"tables,omitempty"`
	Raw    bool              `json:"raw,omitempty"`
}

type DBTablePreview struct {
	Table string   `json:"table"`
	Ops   []string `json:"ops"`
}

type CapabilitiesPreview struct {
	HTTP   *HTTPCapPreview   `json:"http,omitempty"`
	Store  *StoreCapPreview  `json:"store,omitempty"`
	Events *EventsCapPreview `json:"events,omitempty"`
	DB     *DBCapPreview     `json:"db,omitempty"`
}

type WebhookPreview struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

type PipelinePreview struct {
	Name      string `json:"name"`
	Trigger   string `json:"trigger"`
	StepCount int    `json:"step_count"`
}

type DependencyPreview struct {
	ID       string `json:"id"`
	Version  string `json:"version,omitempty"`
	Optional bool   `json:"optional,omitempty"`
}

type PluginPreviewRes struct {
	Name         string              `json:"name"`
	Title        string              `json:"title"`
	Description  string              `json:"description"`
	Version      string              `json:"version"`
	Author       string              `json:"author"`
	Icon         string              `json:"icon"`
	Priority     int                 `json:"priority"`
	HasCSS       bool                `json:"has_css"`
	Capabilities CapabilitiesPreview `json:"capabilities"`
	Depends      []DependencyPreview `json:"depends,omitempty"`
	Settings     []SettingField      `json:"settings"`
	Webhooks     []WebhookPreview    `json:"webhooks"`
	Pipelines    []PipelinePreview   `json:"pipelines"`
	Permissions  []string            `json:"permissions,omitempty"`
}
