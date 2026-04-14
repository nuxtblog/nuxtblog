package sdk

import (
	"context"
	"io"
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
type CategoryDef struct {
	Slug        string `yaml:"slug"          json:"slug"`
	LabelZh     string `yaml:"label_zh"      json:"label_zh"`
	LabelEn     string `yaml:"label_en"      json:"label_en"`
	Order       int    `yaml:"order"         json:"order"`
	MaxPerOwner int    `yaml:"max_per_owner" json:"max_per_owner"`
}
