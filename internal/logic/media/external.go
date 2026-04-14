package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/storage"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"
)

// Link registers an external URL as a media record without uploading a file.
func (s *sMedia) Link(ctx context.Context, req *v1.MediaLinkReq) (*v1.MediaItem, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "upload_files") {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: upload_files")
	}

	uid, _ := middleware.GetCurrentUserID(ctx)

	category := req.Category
	if category == "" {
		category = "post"
	}

	// Derive filename from URL if no title provided.
	filename := req.Title
	if filename == "" {
		parts := strings.Split(strings.TrimRight(req.Url, "/"), "/")
		filename = parts[len(parts)-1]
		if filename == "" {
			filename = "external"
		}
	}

	m := &entity.Medias{
		UploaderId:  int(uid),
		StorageType: int(v1.StorageTypeExternal),
		StorageKey:  "",
		CdnUrl:      req.Url,
		Filename:    filename,
		MimeType:    "image/unknown",
		FileSize:    0,
		AltText:     req.AltText,
		Title:       req.Title,
		Category:    category,
	}

	m.Id = int(idgen.New())
	_, err := dao.Medias.Ctx(ctx).Data(m).Insert()
	if err != nil {
		return nil, fmt.Errorf("save media record: %w", err)
	}
	return s.GetOne(ctx, int64(m.Id))
}

// Localize downloads an external media URL and stores it in local/cloud storage,
// converting the record from StorageTypeExternal to an actual stored file.
func (s *sMedia) Localize(ctx context.Context, req *v1.MediaLocalizeReq) (*v1.MediaItem, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "upload_files") {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: upload_files")
	}

	// Load the existing record.
	var m entity.Medias
	if err := dao.Medias.Ctx(ctx).Where("id", req.Id).Scan(&m); err != nil || m.Id == 0 {
		return nil, errors.New("media not found")
	}
	if v1.StorageType(m.StorageType) != v1.StorageTypeExternal {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "media is not external")
	}

	// Download the external URL.
	resp, err := http.Get(m.CdnUrl) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("download external URL: %w", err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	// Detect MIME type from response header, then sniff from content.
	mimeType := resp.Header.Get("Content-Type")
	if idx := strings.Index(mimeType, ";"); idx > 0 {
		mimeType = strings.TrimSpace(mimeType[:idx])
	}
	if mimeType == "" || mimeType == "application/octet-stream" {
		mimeType = http.DetectContentType(data)
	}

	// Derive file extension.
	ext := ""
	if exts, _ := mime.ExtensionsByType(mimeType); len(exts) > 0 {
		ext = exts[0]
	}
	if ext == "" {
		urlPath := strings.Split(m.CdnUrl, "?")[0]
		parts := strings.Split(urlPath, ".")
		if len(parts) > 1 {
			ext = "." + parts[len(parts)-1]
		}
	}

	// Determine filename.
	filename := m.Filename
	if filename == "" {
		filename = "localized" + ext
	} else if !strings.Contains(filename, ".") {
		filename = filename + ext
	}

	// Resolve storage backend.
	backendName, up := resolveBackend(ctx, mimeType, m.Category)

	// Upload via synthetic multipart file.
	result, err := storage.UploadFromBytes(ctx, up, data, filename, storage.UploadOptions{
		Category:     m.Category,
		PathTemplate: getCategoryPathTemplate(ctx, m.Category),
	})
	if err != nil {
		return nil, fmt.Errorf("upload to storage: %w", err)
	}

	// Determine image dimensions if applicable.
	var width, height int
	if strings.HasPrefix(mimeType, "image/") {
		if img, decErr := imaging.Decode(bytes.NewReader(data)); decErr == nil {
			b := img.Bounds()
			width, height = b.Dx(), b.Dy()
		}
	}

	// Generate thumbnails for images.
	var variantsJSON string
	if strings.HasPrefix(mimeType, "image/") {
		thumbnail, cover, content := getThumbSizes(ctx)
		variants := generateThumbsFromBytes(ctx, up, data, result.StorageKey, thumbnail, cover, content)
		if len(variants) > 0 {
			if b, e := json.Marshal(variants); e == nil {
				variantsJSON = string(b)
			}
		}
	}

	// Update the DB record.
	upd := g.Map{
		"storage_type": result.StorageType,
		"storage_key":  storage.EncodeKey(backendName, result.StorageKey),
		"cdn_url":      result.CdnUrl,
		"mime_type":    mimeType,
		"file_size":    int64(len(data)),
		"variants":     variantsJSON,
	}
	if width > 0 {
		upd["width"] = width
		upd["height"] = height
	}
	if _, err = dao.Medias.Ctx(ctx).Where("id", req.Id).Data(upd).Update(); err != nil {
		_ = up.Delete(ctx, result.StorageKey)
		return nil, fmt.Errorf("update media record: %w", err)
	}

	return s.GetOne(ctx, req.Id)
}
