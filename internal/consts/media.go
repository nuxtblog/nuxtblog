package consts

// MediaCategoryDef declares a media category.
// Built-in entries are synced to options on every startup (upsert, idempotent).
type MediaCategoryDef struct {
	Slug     string // unique slug, used as DB value
	LabelZh  string // Chinese label
	LabelEn  string // English label
	IsSystem bool   // system categories cannot be deleted via API
	Order    int    // display order (ascending)
}

// Media category slug constants — use these throughout the codebase instead of magic strings.
//
// Categories are intentionally broad:
//   - MediaCatUser   covers user avatar AND profile cover; when a user updates either,
//     the previous image for that slot is deleted automatically.
//   - MediaCatPost   covers post featured cover AND inline body images.
//   - MediaCatDoc    covers doc inline images and attachments.
//   - MediaCatMoment covers moment attached images/videos.
//   - MediaCatBanner is for site-wide banner / hero images.
const (
	MediaCatUser   = "user"
	MediaCatPost   = "post"
	MediaCatDoc    = "doc"
	MediaCatMoment = "moment"
	MediaCatBanner = "banner"
)

// BuiltinMediaCategories is the single registration point for all built-in categories.
// To add a new built-in category, append an entry here and restart / run migrate.
var BuiltinMediaCategories = []MediaCategoryDef{
	// 用户 — 头像及个人封面图（更新时会删除旧图）
	{Slug: MediaCatUser, LabelZh: "用户", LabelEn: "User", IsSystem: true, Order: 10},
	// 文章 — 文章封面及正文内嵌图片
	{Slug: MediaCatPost, LabelZh: "文章", LabelEn: "Post", IsSystem: true, Order: 20},
	// 文档 — 文档内嵌图片及附件
	{Slug: MediaCatDoc, LabelZh: "文档", LabelEn: "Doc", IsSystem: true, Order: 30},
	// 动态 — 动态内附图/视频
	{Slug: MediaCatMoment, LabelZh: "动态", LabelEn: "Moment", IsSystem: true, Order: 40},
	// Banner — 站点横幅 / Hero 图
	{Slug: MediaCatBanner, LabelZh: "Banner", LabelEn: "Banner", IsSystem: true, Order: 50},
}

// FileLimitsMB holds per-MIME-class upload size caps in megabytes.
type FileLimitsMB struct {
	Image    float64 `json:"image"`
	Video    float64 `json:"video"`
	Audio    float64 `json:"audio"`
	Document float64 `json:"document"`
	Other    float64 `json:"other"`
}

// DefaultFileLimits is used when the media_size_limits option is absent or zero.
var DefaultFileLimits = FileLimitsMB{Image: 10, Video: 100, Audio: 20, Document: 20, Other: 10}

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
