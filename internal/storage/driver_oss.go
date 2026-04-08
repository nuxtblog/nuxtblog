package storage

// Aliyun OSS storage driver.
// Config keys under storage.backends.<name>:
//
//	type:            "oss"
//	endpoint:        "oss-cn-hangzhou.aliyuncs.com"
//	bucket:          "my-bucket"
//	accessKeyId:     "xxx"
//	accessKeySecret: "xxx"
//	baseUrl:         "https://my-bucket.oss-cn-hangzhou.aliyuncs.com"  # public URL prefix

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

type ossUploader struct {
	bucket  *oss.Bucket
	baseURL string
}

func newOSSUploaderAt(ctx context.Context, prefix string) *ossUploader {
	endpoint, _  := g.Cfg().Get(ctx, prefix+".endpoint")
	bucket, _    := g.Cfg().Get(ctx, prefix+".bucket")
	keyID, _     := g.Cfg().Get(ctx, prefix+".accessKeyId")
	keySecret, _ := g.Cfg().Get(ctx, prefix+".accessKeySecret")
	baseURL, _   := g.Cfg().Get(ctx, prefix+".baseUrl")

	client, err := oss.New(endpoint.String(), keyID.String(), keySecret.String())
	if err != nil {
		g.Log().Warningf(ctx, "[storage/oss] init client failed: %v", err)
		return &ossUploader{baseURL: strings.TrimRight(baseURL.String(), "/")}
	}

	b, err := client.Bucket(bucket.String())
	if err != nil {
		g.Log().Warningf(ctx, "[storage/oss] open bucket failed: %v", err)
		return &ossUploader{baseURL: strings.TrimRight(baseURL.String(), "/")}
	}

	return &ossUploader{
		bucket:  b,
		baseURL: strings.TrimRight(baseURL.String(), "/"),
	}
}

func (u *ossUploader) Upload(_ context.Context, file *ghttp.UploadFile, opts UploadOptions) (*UploadResult, error) {
	if u.bucket == nil {
		return nil, fmt.Errorf("oss bucket not initialised")
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	subDir   := renderPathTemplate(opts.PathTemplate, opts.Category)
	ext      := strings.ToLower(filepath.Ext(file.Filename))
	objKey   := subDir + "/" + randomHex(16) + ext

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	err = u.bucket.PutObject(objKey, src, oss.ContentType(mimeType))
	if err != nil {
		return nil, fmt.Errorf("oss put object: %w", err)
	}

	return &UploadResult{
		StorageType: 3,
		StorageKey:  objKey,
		CdnUrl:      u.baseURL + "/" + objKey,
		MimeType:    mimeType,
		FileSize:    file.Size,
	}, nil
}

func (u *ossUploader) SaveBytes(_ context.Context, data []byte, ext, prefix, originalKey string) (string, error) {
	if u.bucket == nil {
		return "", fmt.Errorf("oss bucket not initialised")
	}

	dir    := filepath.ToSlash(filepath.Dir(originalKey))
	base   := strings.TrimSuffix(filepath.Base(originalKey), filepath.Ext(originalKey))
	objKey := dir + "/" + base + "_" + prefix + ext

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	err := u.bucket.PutObject(objKey, bytes.NewReader(data), oss.ContentType(mimeType))
	if err != nil {
		return "", fmt.Errorf("oss save thumbnail: %w", err)
	}
	return u.baseURL + "/" + objKey, nil
}

func (u *ossUploader) Delete(_ context.Context, storageKey string) error {
	if u.bucket == nil {
		return fmt.Errorf("oss bucket not initialised")
	}
	return u.bucket.DeleteObject(storageKey)
}
