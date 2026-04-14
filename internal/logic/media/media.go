package media

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	v1 "github.com/nuxtblog/nuxtblog/api/media/v1"
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

// ── max_per_owner enforcement ──────────────────────────────────────────────────

// enforceMaxPerOwner deletes the oldest media records for a given uploader+category
// when the count would exceed maxPerOwner after the new upload.
func enforceMaxPerOwner(ctx context.Context, uploaderID int64, category string, maxPerOwner int) {
	var ids []struct {
		Id int `orm:"id"`
	}
	err := dao.Medias.Ctx(ctx).
		Where("uploader_id", uploaderID).
		Where("category", category).
		Where("deleted_at IS NULL").
		OrderAsc("created_at").
		Scan(&ids)
	if err != nil || len(ids) < maxPerOwner {
		return
	}
	// Delete oldest records to make room (keep maxPerOwner-1, since we're about to add one)
	toDelete := ids[:len(ids)-maxPerOwner+1]
	for _, row := range toDelete {
		delUp, actualKey := storage.ForStorageKey(ctx, "")
		var m entity.Medias
		if e := dao.Medias.Ctx(ctx).Where("id", row.Id).Scan(&m); e == nil && m.Id != 0 {
			delUp, actualKey = storage.ForStorageKey(ctx, m.StorageKey)
			_ = delUp.Delete(ctx, actualKey)
		}
		_, _ = dao.Medias.Ctx(ctx).Where("id", row.Id).Delete()
		g.Log().Infof(ctx, "[media] max_per_owner: auto-deleted media %d (uploader=%d, category=%s)", row.Id, uploaderID, category)
	}
}

// ── Upload ────────────────────────────────────────────────────────────────────

func (s *sMedia) Upload(ctx context.Context, req *v1.MediaUploadReq) (*v1.MediaItem, error) {
	role := middleware.GetCurrentUserRole(ctx)
	if !service.Permission().Can(ctx, role, "upload_files") {
		return nil, gerror.NewCode(gcode.CodeNotAuthorized, "permission denied: upload_files")
	}

	mimeType := req.File.Header.Get("Content-Type")
	backendName, up := resolveBackend(ctx, mimeType, req.Category)

	// 1. Validate file extension and size against format policy.
	ext := strings.ToLower(filepath.Ext(req.File.Filename))
	policy := resolveFormatPolicy(ctx, req.Category)
	if err := validateFileExtension(ctx, policy, ext, req.File.Size); err != nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, err.Error())
	}

	// 2. Upload original file.
	result, err := up.Upload(ctx, req.File, storage.UploadOptions{
		Category:     req.Category,
		PathTemplate: getCategoryPathTemplate(ctx, req.Category),
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
		variants := generateThumbsFromFile(ctx, up, req.File, result.StorageKey, thumbnail, cover, content)
		if len(variants) > 0 {
			if b, e := json.Marshal(variants); e == nil {
				variantsJSON = string(b)
			}
		}
	}

	// 4. max_per_owner: auto-delete oldest media for this uploader+category.
	if catDef := GetCategoryDef(ctx, category); catDef != nil && catDef.MaxPerOwner > 0 {
		enforceMaxPerOwner(ctx, uploaderID, category, catDef.MaxPerOwner)
	}

	// 5. Persist media record.
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

// ── Delete ────────────────────────────────────────────────────────────────────

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

// ── Update ────────────────────────────────────────────────────────────────────

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

// ── Query ─────────────────────────────────────────────────────────────────────

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

// ── entity → API item ─────────────────────────────────────────────────────────

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
