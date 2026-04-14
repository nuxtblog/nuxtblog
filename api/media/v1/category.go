package v1

import "github.com/gogf/gf/v2/frame/g"

// MediaCategoryItem is a single category entry returned by the API.
type MediaCategoryItem struct {
	Slug         string `json:"slug"`
	LabelZh      string `json:"label_zh"`
	LabelEn      string `json:"label_en"`
	IsSystem     bool   `json:"is_system"`
	Order        int    `json:"order"`
	StorageKey   string `json:"storage_key"`
	PluginID     string `json:"plugin_id,omitempty"`
	MaxPerOwner  int    `json:"max_per_owner,omitempty"`
	FormatPolicy string `json:"format_policy,omitempty"`
	PathTemplate string `json:"path_template,omitempty"`
}

// ── List ──────────────────────────────────────────────────────────────────────

type MediaCategoryListReq struct {
	g.Meta `path:"/admin/media/categories" method:"get" tags:"Media" summary:"List all media categories"`
}
type MediaCategoryListRes struct {
	List []MediaCategoryItem `json:"list"`
}

// ── Update storage key ────────────────────────────────────────────────────────
// Only storage_key is admin-configurable at runtime.
// Labels and slugs are defined in consts and synced from code on startup.

type MediaCategoryUpdateReq struct {
	g.Meta       `path:"/admin/media/categories/{slug}" method:"put" tags:"Media" summary:"Update category settings"`
	Slug         string  `v:"required" dc:"category slug (path param)"`
	StorageKey   string  `v:"max-length:128" dc:"storage backend name; empty string = use system default"`
	FormatPolicy *string `v:"max-length:64"  dc:"format policy name; empty = use default policy"`
	PathTemplate *string `v:"max-length:128" dc:"path template; empty = use global default"`
}
type MediaCategoryUpdateRes struct{}
