package consts

// MediaCategoryDef declares a media category.
// Built-in entries are synced to options on every startup (upsert, idempotent).
type MediaCategoryDef struct {
	Slug        string // unique slug, used as DB value
	LabelZh     string // Chinese label
	LabelEn     string // English label
	IsSystem    bool   // system categories cannot be deleted via API
	Order       int    // display order (ascending)
	MaxPerOwner int    // >0 = auto-replace oldest when uploader exceeds this count; 0 = unlimited
}

// Media category slug constants — use these throughout the codebase instead of magic strings.
const (
	MediaCatAvatar = "avatar"
	MediaCatCover  = "cover"
	MediaCatPost   = "post"
	MediaCatDoc    = "doc"
	MediaCatMoment = "moment"
	MediaCatBanner = "banner"
)

// BuiltinMediaCategories is the single registration point for all built-in categories.
// To add a new built-in category, append an entry here and restart / run migrate.
var BuiltinMediaCategories = []MediaCategoryDef{
	{Slug: MediaCatAvatar, LabelZh: "头像", LabelEn: "Avatar", IsSystem: true, Order: 5, MaxPerOwner: 1},
	{Slug: MediaCatCover, LabelZh: "封面", LabelEn: "Cover", IsSystem: true, Order: 6, MaxPerOwner: 1},
	{Slug: MediaCatPost, LabelZh: "文章", LabelEn: "Post", IsSystem: true, Order: 20, MaxPerOwner: 0},
	{Slug: MediaCatDoc, LabelZh: "文档", LabelEn: "Doc", IsSystem: true, Order: 30, MaxPerOwner: 0},
	{Slug: MediaCatMoment, LabelZh: "动态", LabelEn: "Moment", IsSystem: true, Order: 40, MaxPerOwner: 0},
	{Slug: MediaCatBanner, LabelZh: "横幅", LabelEn: "Banner", IsSystem: true, Order: 50, MaxPerOwner: 0},
}

// ExtensionGroup defines a set of file extensions and their size limit.
type ExtensionGroup struct {
	Name       string   `json:"name"`
	LabelZh    string   `json:"label_zh"`
	LabelEn    string   `json:"label_en"`
	Extensions []string `json:"extensions"`
	MaxSizeMB  float64  `json:"max_size_mb"`
}

// FormatPolicy combines multiple ExtensionGroups and is bound to categories.
type FormatPolicy struct {
	Name     string   `json:"name"`
	LabelZh  string   `json:"label_zh"`
	LabelEn  string   `json:"label_en"`
	IsSystem bool     `json:"is_system"`
	Groups   []string `json:"groups"`
}

// DefaultExtensionGroups are the built-in extension group definitions.
var DefaultExtensionGroups = []ExtensionGroup{
	{Name: "images", LabelZh: "图片", LabelEn: "Images", Extensions: []string{"jpg", "jpeg", "png", "webp", "gif", "svg", "ico", "bmp"}, MaxSizeMB: 10},
	{Name: "videos", LabelZh: "视频", LabelEn: "Videos", Extensions: []string{"mp4", "webm", "mov", "avi", "mkv"}, MaxSizeMB: 100},
	{Name: "audio", LabelZh: "音频", LabelEn: "Audio", Extensions: []string{"mp3", "wav", "ogg", "flac", "aac", "m4a"}, MaxSizeMB: 20},
	{Name: "documents", LabelZh: "文档", LabelEn: "Documents", Extensions: []string{"pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "md"}, MaxSizeMB: 20},
}

// DefaultFormatPolicies are the built-in format policy definitions.
var DefaultFormatPolicies = []FormatPolicy{
	{Name: "default", LabelZh: "默认策略", LabelEn: "Default", IsSystem: true, Groups: []string{"images", "videos", "audio", "documents"}},
	{Name: "images_only", LabelZh: "仅图片", LabelEn: "Images Only", IsSystem: true, Groups: []string{"images"}},
}

// ThumbSize describes a thumbnail preset dimension. Height=0 means "auto (keep ratio)".
type ThumbSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// DefaultThumb* are used when the corresponding thumbnail options are absent.
var (
	DefaultThumbThumbnail = ThumbSize{Width: 300, Height: 200}
	DefaultThumbCover     = ThumbSize{Width: 400, Height: 300}
	DefaultThumbContent   = ThumbSize{Width: 1200, Height: 0}
)
