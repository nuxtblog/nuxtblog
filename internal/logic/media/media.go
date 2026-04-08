package media

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"io"
	"mime"
	"net/http"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/dao"
	"github.com/nuxtblog/nuxtblog/internal/event"
	"github.com/nuxtblog/nuxtblog/internal/event/payload"
	"github.com/nuxtblog/nuxtblog/internal/middleware"
	"github.com/nuxtblog/nuxtblog/internal/model/entity"
	eng "github.com/nuxtblog/nuxtblog/internal/pluginsys"
	"github.com/nuxtblog/nuxtblog/internal/service"
	"github.com/nuxtblog/nuxtblog/internal/storage"
	"github.com/nuxtblog/nuxtblog/internal/util/idgen"
)

type sMedia struct{}

func New() service.IMedia { return &sMedia{} }

func init() {
	service.RegisterMedia(New())
}

// ── option helpers ────────────────────────────────────────────────────────────

func getOptionJSON(ctx context.Context, key string, out any) error {
	val, err := dao.Options.Ctx(ctx).Where("key", key).Value("value")
	if err != nil || val.IsNil() {
		return err
	}
	raw := val.String()
	// DB value is JSON-encoded string wrapping a JSON object, e.g. `"{\"image\":10}"`
	var inner string
	if json.Unmarshal([]byte(raw), &inner) == nil {
		return json.Unmarshal([]byte(inner), out)
	}
	return json.Unmarshal([]byte(raw), out)
}

// ── size limits ───────────────────────────────────────────────────────────────

func getLimits(ctx context.Context) consts.FileLimitsMB {
	var lim consts.FileLimitsMB
	if err := getOptionJSON(ctx, "media_size_limits", &lim); err != nil || lim.Image == 0 {
		return consts.DefaultFileLimits
	}
	return lim
}

func limitMBForMime(lim consts.FileLimitsMB, mimeType string) float64 {
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return lim.Image
	case strings.HasPrefix(mimeType, "video/"):
		return lim.Video
	case strings.HasPrefix(mimeType, "audio/"):
		return lim.Audio
	case strings.HasPrefix(mimeType, "application/pdf"),
		strings.HasPrefix(mimeType, "application/msword"),
		strings.HasPrefix(mimeType, "application/vnd.openxmlformats"):
		return lim.Document
	default:
		return lim.Other
	}
}

// ── thumbnail config ──────────────────────────────────────────────────────────

func getThumbSizes(ctx context.Context) (thumbnail, cover, content consts.ThumbSize) {
	thumbnail = consts.DefaultThumbThumbnail
	cover     = consts.DefaultThumbCover
	content   = consts.DefaultThumbContent
	_ = getOptionJSON(ctx, "media_thumbnail", &thumbnail)
	_ = getOptionJSON(ctx, "media_cover_thumb", &cover)
	_ = getOptionJSON(ctx, "media_content_thumb", &content)
	return
}

// generateThumbs reads image bytes from the upload file, resizes/crops them,
// and saves thumbnails via the ThumbSaver interface (local storage).
// Errors are non-fatal: logged and skipped.
func generateThumbs(ctx context.Context, up storage.Uploader, file *ghttp.UploadFile, originalKey string, thumbnail, cover, content consts.ThumbSize) map[string]string {
	ts, ok := up.(storage.ThumbSaver)
	if !ok {
		return nil
	}

	// Re-open the multipart file to get image bytes.
	src, err := file.Open()
	if err != nil {
		g.Log().Warningf(ctx, "thumbnail: open file: %v", err)
		return nil
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		g.Log().Warningf(ctx, "thumbnail: read file: %v", err)
		return nil
	}

	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		g.Log().Warningf(ctx, "thumbnail: decode image: %v", err)
		return nil
	}

	variants := make(map[string]string)

	save := func(name string, w, h int) {
		if w <= 0 && h <= 0 {
			return
		}
		var resized = img
		if w > 0 && h > 0 {
			// both set → crop to exact size
			resized = imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)
		} else {
			// one of w/h is 0 → imaging treats 0 as "auto, keep ratio"
			resized = imaging.Resize(img, w, h, imaging.Lanczos)
		}
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 85}); err != nil {
			g.Log().Warningf(ctx, "thumbnail: encode %s: %v", name, err)
			return
		}
		cdnUrl, err := ts.SaveBytes(ctx, buf.Bytes(), ".jpg", name, originalKey)
		if err != nil {
			g.Log().Warningf(ctx, "thumbnail: save %s: %v", name, err)
			return
		}
		variants[name] = cdnUrl
	}

	save("thumbnail", thumbnail.Width, thumbnail.Height)
	save("cover", cover.Width, cover.Height)
	save("content", content.Width, content.Height)

	return variants
}

// ── Upload path ───────────────────────────────────────────────────────────────

func getUploadPathTemplate(ctx context.Context) string {
	val, err := dao.Options.Ctx(ctx).Where("key", "media_upload_path").Value("value")
	if err != nil || val.IsNil() {
		return consts.StoragePathTemplateYearMonth
	}
	raw := val.String()
	var s string
	if json.Unmarshal([]byte(raw), &s) == nil && s != "" {
		return s
	}
	if raw != "" {
		return raw
	}
	return consts.StoragePathTemplateYearMonth
}

// ── Storage routing ───────────────────────────────────────────────────────────

type storageRule struct {
	MimePrefix string `json:"mimePrefix"`
	Backend    string `json:"backend"`
}

type storageRouting struct {
	Default string        `json:"default"`
	Rules   []storageRule `json:"rules"`
}

// resolveBackend picks the right storage backend.
// Priority: category storage key (from media_categories option) > MIME prefix rule > default.
func resolveBackend(ctx context.Context, mimeType, category string) (backendName string, up storage.Uploader) {
	// 1. Category storage key from media_categories option (takes highest priority)
	if category != "" {
		if key := GetCategoryStorageKey(ctx, category); key != "" {
			return key, storage.Named(ctx, key)
		}
	}

	// 2. MIME prefix rules from storage_routing option
	var routing storageRouting
	_ = getOptionJSON(ctx, "storage_routing", &routing)
	for _, rule := range routing.Rules {
		if rule.MimePrefix != "" && strings.HasPrefix(mimeType, rule.MimePrefix) {
			return rule.Backend, storage.Named(ctx, rule.Backend)
		}
	}

	// 3. Routing default, then config default
	if routing.Default != "" {
		return routing.Default, storage.Named(ctx, routing.Default)
	}
	name := storage.DefaultName(ctx)
	return name, storage.Named(ctx, name)
}

// ── Upload ────────────────────────────────────────────────────────────────────

func (s *sMedia) Upload(ctx context.Context, req *v1.MediaUploadReq) (*v1.MediaItem, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "upload_files") {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: upload_files")
	}

	mimeType := req.File.Header.Get("Content-Type")
	backendName, up := resolveBackend(ctx, mimeType, req.Category)

	// 1. Validate file size against configured limits.
	lim := getLimits(ctx)
	maxBytes := int64(limitMBForMime(lim, mimeType) * 1024 * 1024)
	if req.File.Size > maxBytes {
		mb := float64(req.File.Size) / 1024 / 1024
		maxMB := limitMBForMime(lim, mimeType)
		return nil, gerror.NewCode(gcode.CodeInvalidParameter,
			g.I18n().Tf(ctx, "media.file_too_large", mb, maxMB))
	}

	// 2. Upload original file.
	result, err := up.Upload(ctx, req.File, storage.UploadOptions{
		Category:     req.Category,
		PathTemplate: getUploadPathTemplate(ctx),
	})
	if err != nil {
		return nil, fmt.Errorf("upload file: %w", err)
	}

	uploaderID, _ := middleware.GetCurrentUserID(ctx)

	// Run plugin filter:media.upload — allows plugins to modify metadata or reject uploads
	{
		altText := req.AltText
		title := req.Title
		filterIn := map[string]any{
			"filename":  req.File.Filename,
			"mime_type": result.MimeType,
			"category":  req.Category,
			"alt_text":  altText,
			"title":     title,
		}
		if filtered, ferr := eng.Filter(ctx, eng.FilterMediaUpload, filterIn); ferr != nil {
			return nil, ferr
		} else {
			if v, ok := filtered["alt_text"].(string); ok {
				req.AltText = v
			}
			if v, ok := filtered["title"].(string); ok {
				req.Title = v
			}
		}
	}

	category := req.Category
	if category == "" {
		category = "post"
	}

	// 3. Generate thumbnails for images (best-effort, errors are logged not returned).
	var variantsJSON string
	if strings.HasPrefix(result.MimeType, "image/") {
		thumbnail, cover, content := getThumbSizes(ctx)
		variants := generateThumbs(ctx, up, req.File, result.StorageKey, thumbnail, cover, content)
		if len(variants) > 0 {
			if b, e := json.Marshal(variants); e == nil {
				variantsJSON = string(b)
			}
		}
	}

	// 4. Persist media record.
	m := &entity.Medias{
		UploaderId:  int(uploaderID),
		StorageType: result.StorageType,
		StorageKey:  storage.EncodeKey(backendName, result.StorageKey),
		CdnUrl:      result.CdnUrl,
		Filename:    req.File.Filename,
		MimeType:    result.MimeType,
		FileSize:    int(result.FileSize),
		Width:       result.Width,
		Height:      result.Height,
		Duration:    result.Duration,
		AltText:     req.AltText,
		Title:       req.Title,
		Category:    category,
		Variants:    variantsJSON,
	}

	// Assign a random non-sequential ID to avoid enumeration / cleanup confusion.
	m.Id = int(idgen.New())
	_, err = dao.Medias.Ctx(ctx).Data(m).Insert()
	if err != nil {
		_ = up.Delete(ctx, result.StorageKey)
		return nil, fmt.Errorf("save media record: %w", err)
	}
	id := int64(m.Id)

	item, err := s.GetOne(ctx, id)
	if err != nil {
		return nil, err
	}
	_ = event.Emit(ctx, event.MediaUploaded, payload.MediaUploaded{
		MediaID:    id,
		UploaderID: uploaderID,
		Filename:   req.File.Filename,
		MimeType:   result.MimeType,
		FileSize:   result.FileSize,
		URL:        item.CdnUrl,
		Category:   category,
		Width:      result.Width,
		Height:     result.Height,
	})
	return item, nil
}

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

	// Assign a random non-sequential ID (same reason as Upload).
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
		// Fall back: try to extract from URL.
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
		PathTemplate: getUploadPathTemplate(ctx),
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
		if ts, ok := up.(storage.ThumbSaver); ok {
			thumbnail, cover, content := getThumbSizes(ctx)
			img, decErr := imaging.Decode(bytes.NewReader(data))
			if decErr == nil {
				variants := make(map[string]string)
				saveFn := func(name string, w, h int) {
					if w <= 0 && h <= 0 {
						return
					}
					var resized = img
					if w > 0 && h > 0 {
						resized = imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)
					} else {
						resized = imaging.Resize(img, w, h, imaging.Lanczos)
					}
					var buf bytes.Buffer
					if encErr := jpeg.Encode(&buf, resized, &jpeg.Options{Quality: 85}); encErr != nil {
						g.Log().Warningf(ctx, "localize thumbnail: encode %s: %v", name, encErr)
						return
					}
					cdnUrl, saveErr := ts.SaveBytes(ctx, buf.Bytes(), ".jpg", name, result.StorageKey)
					if saveErr != nil {
						g.Log().Warningf(ctx, "localize thumbnail: save %s: %v", name, saveErr)
						return
					}
					variants[name] = cdnUrl
				}
				saveFn("thumbnail", thumbnail.Width, thumbnail.Height)
				saveFn("cover", cover.Width, cover.Height)
				saveFn("content", content.Width, content.Height)
				if len(variants) > 0 {
					if b, e := json.Marshal(variants); e == nil {
						variantsJSON = string(b)
					}
				}
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

func (s *sMedia) Delete(ctx context.Context, id int64) error {
	var m entity.Medias
	err := dao.Medias.Ctx(ctx).Where("id", id).Scan(&m)
	if err != nil || m.Id == 0 {
		return errors.New(g.I18n().T(ctx, "error.media_not_found"))
	}

	// For admin-level users, verify delete_others_media capability.
	role := middleware.GetCurrentUserRole(ctx)
	if role >= middleware.RoleAdmin {
		uid, _ := middleware.GetCurrentUserID(ctx)
		if int64(m.UploaderId) != uid && !service.Permission().Can(ctx, role, "delete_others_media") {
			return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: delete_others_media")
		}
	}

	delUp, actualKey := storage.ForStorageKey(ctx, m.StorageKey)
	if err = delUp.Delete(ctx, actualKey); err != nil {
		g.Log().Warningf(ctx, "delete storage file: %v", err)
	}

	uploaderID, _ := middleware.GetCurrentUserID(ctx)
	_, err = dao.Medias.Ctx(ctx).Where("id", id).Delete()
	if err != nil {
		return err
	}
	_ = event.Emit(ctx, event.MediaDeleted, payload.MediaDeleted{
		MediaID:    id,
		UploaderID: uploaderID,
		Filename:   m.Filename,
		MimeType:   m.MimeType,
		Category:   m.Category,
	})
	return nil
}

func (s *sMedia) Update(ctx context.Context, req *v1.MediaUpdateReq) error {
	// For admin-level users, verify manage_media capability when editing others' files.
	role := middleware.GetCurrentUserRole(ctx)
	if role >= middleware.RoleAdmin {
		uid, _ := middleware.GetCurrentUserID(ctx)
		type ownerRow struct{ UploaderId int `orm:"uploader_id"` }
		var row ownerRow
		if e := dao.Medias.Ctx(ctx).Where("id", req.Id).Scan(&row); e == nil && int64(row.UploaderId) != uid {
			if !service.Permission().Can(ctx, role, "manage_media") {
				return gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: manage_media")
			}
		}
	}

	data := g.Map{}
	if req.AltText != nil {
		data["alt_text"] = *req.AltText
	}
	if req.Title != nil {
		data["title"] = *req.Title
	}
	if len(data) == 0 {
		return nil
	}
	_, err := dao.Medias.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	return err
}

func (s *sMedia) GetOne(ctx context.Context, id int64) (*v1.MediaItem, error) {
	var m entity.Medias
	err := dao.Medias.Ctx(ctx).Where("id", id).Scan(&m)
	if err != nil || m.Id == 0 {
		return nil, errors.New(g.I18n().T(ctx, "error.media_not_found"))
	}
	return entityToItem(&m), nil
}

func (s *sMedia) GetList(ctx context.Context, req *v1.MediaGetListReq) (*v1.MediaGetListRes, error) {
	model := dao.Medias.Ctx(ctx)
	if req.MimeType != nil {
		model = model.WhereLike("mime_type", *req.MimeType+"%")
	}
	if req.UploaderId != nil {
		model = model.Where("uploader_id", *req.UploaderId)
	}
	if req.Category != nil && *req.Category != "" {
		model = model.Where("category", *req.Category)
	}
	if req.StorageType != nil {
		model = model.Where("storage_type", *req.StorageType)
	}

	total, err := model.Count()
	if err != nil {
		return nil, err
	}

	var list []entity.Medias
	err = model.OrderDesc("created_at").Page(req.Page, req.Size).Scan(&list)
	if err != nil {
		return nil, err
	}

	items := make([]*v1.MediaItem, 0, len(list))
	for i := range list {
		items = append(items, entityToItem(&list[i]))
	}

	return &v1.MediaGetListRes{
		List:  items,
		Total: total,
		Page:  req.Page,
		Size:  req.Size,
	}, nil
}

func (s *sMedia) GetStats(ctx context.Context) (*v1.MediaGetStatsRes, error) {
	type row struct {
		MimePrefix string `orm:"mime_prefix"`
		Count      int    `orm:"count"`
	}
	var rows []row
	err := dao.Medias.DB().Ctx(ctx).Raw(`
		SELECT
			CASE
				WHEN mime_type LIKE 'image/%' THEN 'image'
				WHEN mime_type LIKE 'video/%' THEN 'video'
				WHEN mime_type LIKE 'audio/%' THEN 'audio'
				ELSE 'other'
			END AS mime_prefix,
			COUNT(*) AS count
		FROM medias
		WHERE deleted_at IS NULL
		GROUP BY mime_prefix`).Scan(&rows)
	if err != nil {
		return nil, err
	}
	res := &v1.MediaGetStatsRes{}
	for _, r := range rows {
		res.Total += r.Count
		switch r.MimePrefix {
		case "image":
			res.Image = r.Count
		case "video":
			res.Video = r.Count
		case "audio":
			res.Audio = r.Count
		case "other":
			res.Other = r.Count
		}
	}
	return res, nil
}

func entityToItem(m *entity.Medias) *v1.MediaItem {
	item := &v1.MediaItem{
		Id:          int64(m.Id),
		UploaderId:  int64(m.UploaderId),
		StorageType: v1.StorageType(m.StorageType),
		CdnUrl:      m.CdnUrl,
		Filename:    m.Filename,
		MimeType:    m.MimeType,
		FileSize:    int64(m.FileSize),
		AltText:     m.AltText,
		Title:       m.Title,
		Category:    m.Category,
		Variants:    m.Variants,
		CreatedAt:   m.CreatedAt,
	}
	if m.Width > 0 {
		w := m.Width
		item.Width = &w
	}
	if m.Height > 0 {
		h := m.Height
		item.Height = &h
	}
	if m.Duration > 0 {
		d := m.Duration
		item.Duration = &d
	}
	return item
}
