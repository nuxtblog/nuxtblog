package sdk

import (
	"context"
	"io"
	"strings"
)

// StorageAdapter is the interface that cloud storage plugins implement
// to provide upload/delete capabilities.
type StorageAdapter interface {
	Upload(ctx context.Context, file *UploadFile, opts UploadOpts) (*UploadResult, error)
	Delete(ctx context.Context, storageKey string) error
}

// ThumbAdapter is optionally implemented by StorageAdapter to support
// saving pre-generated thumbnail bytes alongside the original file.
type ThumbAdapter interface {
	SaveBytes(ctx context.Context, data []byte, ext, prefix, originalKey string) (cdnUrl string, err error)
}

// UploadFile carries file data for plugin-initiated uploads.
type UploadFile struct {
	Filename string
	Size     int64
	Reader   io.ReadSeeker
}

// UploadOpts carries per-upload metadata.
type UploadOpts struct {
	Category     string
	PathTemplate string
}

// UploadResult holds the outcome of a successful upload.
type UploadResult struct {
	StorageType int
	StorageKey  string
	CdnUrl      string
	MimeType    string
	FileSize    int64
	Width       int
	Height      int
	Duration    int
	Variants    map[string]string
}

// MediaService is the interface exposed to plugins via PluginContext.Media.
type MediaService interface {
	RegisterStorageAdapter(name string, storageType int, adapter StorageAdapter) error
	RegisterCategory(def CategoryDef) error
	Upload(ctx context.Context, data []byte, filename string, opts UploadOpts) (*UploadResult, error)
	Delete(ctx context.Context, mediaID int64) error
}

// CategoryDef declares a media category that a plugin can register.
// Label supports plain strings or {{key}} templates resolved via the plugin's i18n block.
type CategoryDef struct {
	Slug        string `yaml:"slug"          json:"slug"`
	Label       string `yaml:"label"         json:"label"`
	Order       int    `yaml:"order"         json:"order"`
	MaxPerOwner int    `yaml:"max_per_owner" json:"max_per_owner"`

	// ResolvedZh/ResolvedEn are populated by ResolveCategoryLabel for internal use.
	// Not serialized from plugin.yaml.
	ResolvedZh string `yaml:"-" json:"-"`
	ResolvedEn string `yaml:"-" json:"-"`
}

// ResolveCategoryLabel resolves Label's {{key}} template using the plugin's
// i18n block, populating ResolvedZh/ResolvedEn for internal storage.
// If Label contains no {{}} templates, both locales get the plain string.
func (d *CategoryDef) ResolveCategoryLabel(i18n map[string]map[string]string) {
	if d.Label == "" {
		return
	}
	resolve := func(locale string) string {
		msgs := i18n[locale]
		if msgs == nil {
			return d.Label
		}
		result := d.Label
		for {
			start := strings.Index(result, "{{")
			if start < 0 {
				break
			}
			end := strings.Index(result[start:], "}}")
			if end < 0 {
				break
			}
			key := result[start+2 : start+end]
			val, ok := msgs[key]
			if !ok {
				break
			}
			result = result[:start] + val + result[start+end+2:]
		}
		return result
	}
	d.ResolvedZh = resolve("zh")
	d.ResolvedEn = resolve("en")
	// If en is still a template (no en translation), fall back to zh
	if d.ResolvedEn == d.Label && d.ResolvedZh != d.Label {
		d.ResolvedEn = d.ResolvedZh
	}
}

