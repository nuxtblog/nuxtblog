package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/internal/consts"
)

// UploadResult holds the result of a successful upload.
type UploadResult struct {
	StorageType int               // 1=local 2=S3/MinIO 3=OSS 4=COS
	StorageKey  string            // internal key used for deletion
	CdnUrl      string            // public access URL (always absolute)
	MimeType    string
	FileSize    int64
	Width       int               // 0 if not an image
	Height      int               // 0 if not an image
	Duration    int               // 0 if not a video/audio
	Variants    map[string]string // thumbnail URLs keyed by name ("thumbnail", "cover", "content")
}

// UploadOptions carries per-upload metadata passed through to storage drivers.
type UploadOptions struct {
	Category     string // media category slug (e.g. "post_content", "avatar")
	PathTemplate string // path template override; "" = use default "{year}/{month}"
}

// Uploader is the pluggable interface for file storage backends.
type Uploader interface {
	Upload(ctx context.Context, file *ghttp.UploadFile, opts UploadOptions) (*UploadResult, error)
	Delete(ctx context.Context, storageKey string) error
}

// ThumbSaver is an optional interface implemented by backends that support
// saving pre-generated thumbnail bytes alongside the original file.
type ThumbSaver interface {
	SaveBytes(ctx context.Context, data []byte, ext, prefix, originalKey string) (cdnUrl string, err error)
}

// BackendInfo describes a configured storage backend.
type BackendInfo struct {
	Name        string `json:"name"`         // config key, used as ID
	DisplayName string `json:"display_name"` // human-readable label from config
	Type        string `json:"type"`
	Enabled     bool   `json:"enabled"`
}

// ── Public API ────────────────────────────────────────────────────────────────

// Named returns the Uploader for a backend defined in config under
// storage.backends.<name>.  Returns the local fallback if the backend is
// missing or disabled.
func Named(ctx context.Context, name string) Uploader {
	return buildBackend(ctx, name)
}

// DefaultName returns the name of the default backend.
func DefaultName(ctx context.Context) string {
	return defaultBackendName(ctx)
}

// Default returns the Uploader for the configured default backend.
// Supports both new (storage.default) and legacy (storage.driver) config.
func Default(ctx context.Context) Uploader {
	return Named(ctx, defaultBackendName(ctx))
}

// ForStorageKey parses an encoded storage key ("backendName|actualKey") and
// returns the appropriate Uploader and the actual key to pass to Delete.
// Legacy keys without a backend prefix fall back to Default().
func ForStorageKey(ctx context.Context, encoded string) (up Uploader, actualKey string) {
	name, key := parseKey(encoded)
	if name == "" {
		return Default(ctx), key
	}
	return Named(ctx, name), key
}

// EncodeKey produces the DB-storable composite key: "backendName|actualKey".
func EncodeKey(backendName, actualKey string) string {
	return backendName + "|" + actualKey
}

// ListBackends returns all backends declared in storage.backends, or a
// synthesised single-entry list for legacy config.
func ListBackends(ctx context.Context) []BackendInfo {
	backendsVal, err := g.Cfg().Get(ctx, "storage.backends")
	if err != nil || backendsVal.IsNil() {
		return legacySynthesise(ctx)
	}
	m := backendsVal.Map()
	infos := make([]BackendInfo, 0, len(m))
	for name := range m {
		prefix := "storage.backends." + name
		typeVal, _ := g.Cfg().Get(ctx, prefix+".type")
		enabledVal, _ := g.Cfg().Get(ctx, prefix+".enabled")
		displayVal, _ := g.Cfg().Get(ctx, prefix+".name")
		enabled := enabledVal.IsNil() || enabledVal.Bool()
		displayName := displayVal.String()
		if displayName == "" {
			displayName = name
		}
		infos = append(infos, BackendInfo{
			Name:        name,
			DisplayName: displayName,
			Type:        typeVal.String(),
			Enabled:     enabled,
		})
	}
	return infos
}

// ── Internal ──────────────────────────────────────────────────────────────────

func defaultBackendName(ctx context.Context) string {
	if v, _ := g.Cfg().Get(ctx, "storage.default"); v.String() != "" {
		return v.String()
	}
	return consts.StorageDriverLocal
}

func buildBackend(ctx context.Context, name string) Uploader {
	prefix := "storage.backends." + name
	typeVal, _ := g.Cfg().Get(ctx, prefix+".type")
	enabledVal, _ := g.Cfg().Get(ctx, prefix+".enabled")

	// enabled defaults to true when the key is absent
	if !enabledVal.IsNil() && !enabledVal.Bool() {
		g.Log().Warningf(ctx, "[storage] backend %q is disabled, using local fallback", name)
		return newLocalUploaderAt(ctx, "storage.local")
	}

	switch typeVal.String() {
	case consts.StorageDriverLocal:
		return newLocalUploaderAt(ctx, prefix)
	case consts.StorageDriverS3:
		return newS3UploaderAt(ctx, prefix)
	case consts.StorageDriverOSS:
		return newOSSUploaderAt(ctx, prefix)
	case consts.StorageDriverCOS:
		return newCOSUploaderAt(ctx, prefix)
	default:
		// No "type" key → legacy single-backend config
		return buildLegacy(ctx, name)
	}
}

func buildLegacy(ctx context.Context, _ string) Uploader {
	return newLocalUploaderAt(ctx, "storage.local")
}

func legacySynthesise(_ context.Context) []BackendInfo {
	return []BackendInfo{{Name: consts.StorageDriverLocal, Type: consts.StorageDriverLocal, Enabled: true}}
}

// parseKey splits "name|key" → ("name", "key"), or "key" → ("", "key").
func parseKey(s string) (name, key string) {
	if i := strings.IndexByte(s, '|'); i >= 0 {
		return s[:i], s[i+1:]
	}
	return "", s
}

// UploadFromBytes uploads raw bytes through the given Uploader by creating a
// synthetic multipart.FileHeader backed by an in-memory buffer.
// filename should include the file extension (e.g. "image.jpg").
func UploadFromBytes(ctx context.Context, up Uploader, data []byte, filename string, opts UploadOptions) (*UploadResult, error) {
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		return nil, fmt.Errorf("storage.UploadFromBytes: create form file: %w", err)
	}
	if _, err = io.Copy(fw, bytes.NewReader(data)); err != nil {
		return nil, fmt.Errorf("storage.UploadFromBytes: copy data: %w", err)
	}
	mw.Close()

	mr := multipart.NewReader(&buf, mw.Boundary())
	form, err := mr.ReadForm(int64(len(data)) + 4096)
	if err != nil {
		return nil, fmt.Errorf("storage.UploadFromBytes: read form: %w", err)
	}
	defer form.RemoveAll()

	fileHeaders := form.File["file"]
	if len(fileHeaders) == 0 {
		return nil, fmt.Errorf("storage.UploadFromBytes: no file header parsed")
	}

	uf := &ghttp.UploadFile{FileHeader: fileHeaders[0]}
	return up.Upload(ctx, uf, opts)
}
