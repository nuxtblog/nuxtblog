package storage

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/internal/consts"
)

type localUploader struct {
	uploadPath string
	baseURL    string
}

func newLocalUploaderAt(ctx context.Context, prefix string) *localUploader {
	path, _ := g.Cfg().Get(ctx, prefix+".uploadPath")
	base, _ := g.Cfg().Get(ctx, prefix+".baseUrl")
	uploadPath := path.String()
	if uploadPath == "" {
		uploadPath = "./uploads"
	}
	baseURL := strings.TrimRight(base.String(), "/")
	if baseURL == "" {
		baseURL = "http://localhost:9000/uploads"
	}
	return &localUploader{uploadPath: uploadPath, baseURL: baseURL}
}

func renderPathTemplate(tmpl, category string) string {
	if tmpl == "" {
		tmpl = consts.StoragePathTemplateYearMonth
	}
	now := time.Now()
	// sanitise category so it is safe as a path segment
	var catSafe strings.Builder
	for _, r := range category {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' {
			catSafe.WriteRune(r)
		}
	}
	replacer := strings.NewReplacer(
		consts.StoragePathPlaceholderYear,     fmt.Sprintf("%d", now.Year()),
		consts.StoragePathPlaceholderMonth,    fmt.Sprintf("%02d", int(now.Month())),
		consts.StoragePathPlaceholderDay,      fmt.Sprintf("%02d", now.Day()),
		consts.StoragePathPlaceholderCategory, catSafe.String(),
	)
	result := replacer.Replace(tmpl)
	// collapse consecutive slashes and strip leading/trailing slashes
	parts := strings.Split(result, "/")
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}
	if len(cleaned) == 0 {
		return fmt.Sprintf("%d/%02d", now.Year(), int(now.Month()))
	}
	return strings.Join(cleaned, "/")
}

func (u *localUploader) Upload(_ context.Context, file *ghttp.UploadFile, opts UploadOptions) (*UploadResult, error) {
	subDir := renderPathTemplate(opts.PathTemplate, opts.Category)
	dir := filepath.Join(u.uploadPath, filepath.FromSlash(subDir))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := randomHex(16) + ext
	savedPath := filepath.Join(dir, filename)

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open upload file: %w", err)
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, fmt.Errorf("read upload file: %w", err)
	}
	if err := os.WriteFile(savedPath, data, 0644); err != nil {
		return nil, fmt.Errorf("save file: %w", err)
	}

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	var width, height int
	if strings.HasPrefix(mimeType, "image/") {
		width, height = readImageDimensions(savedPath)
	}

	storageKey := subDir + "/" + filename
	return &UploadResult{
		StorageType: 1,
		StorageKey:  storageKey,
		CdnUrl:      u.baseURL + "/" + storageKey,
		MimeType:    mimeType,
		FileSize:    int64(len(data)),
		Width:       width,
		Height:      height,
	}, nil
}

func (u *localUploader) SaveBytes(_ context.Context, data []byte, ext, prefix, originalKey string) (string, error) {
	dir := filepath.Join(u.uploadPath, filepath.Dir(filepath.FromSlash(originalKey)))
	base := strings.TrimSuffix(filepath.Base(filepath.FromSlash(originalKey)), filepath.Ext(filepath.FromSlash(originalKey)))
	filename := base + "_" + prefix + ext
	fullPath := filepath.Join(dir, filename)
	if err := os.WriteFile(fullPath, data, 0644); err != nil {
		return "", fmt.Errorf("save thumbnail: %w", err)
	}
	subDir := filepath.ToSlash(filepath.Dir(originalKey))
	return u.baseURL + "/" + subDir + "/" + filename, nil
}

func (u *localUploader) Delete(_ context.Context, storageKey string) error {
	path := filepath.Join(u.uploadPath, filepath.FromSlash(storageKey))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(path)
}

func readImageDimensions(path string) (width, height int) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	cfg, _, err := image.DecodeConfig(f)
	if err != nil {
		return
	}
	return cfg.Width, cfg.Height
}

func randomHex(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return hex.EncodeToString(b)
}
