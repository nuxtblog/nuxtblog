package storage

// S3-compatible storage driver (AWS S3, MinIO, Cloudflare R2, …)
// Config keys under storage.backends.<name>:
//
//	type:      "s3"
//	endpoint:  "s3.amazonaws.com"          # or MinIO host:port / R2 endpoint
//	bucket:    "my-bucket"
//	region:    "us-east-1"
//	accessKey: "AKIAIOSFODNN7EXAMPLE"
//	secretKey: "wJalrXUtnFEMI..."
//	baseUrl:   "https://my-bucket.s3.amazonaws.com"   # public URL prefix
//	useSSL:    true                         # default true
//	pathStyle: false                        # true for self-hosted MinIO

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type s3Uploader struct {
	client    *minio.Client
	bucket    string
	baseURL   string
	uploadPath string // base prefix inside bucket (optional)
}

func newS3UploaderAt(ctx context.Context, prefix string) *s3Uploader {
	endpoint, _  := g.Cfg().Get(ctx, prefix+".endpoint")
	bucket, _    := g.Cfg().Get(ctx, prefix+".bucket")
	region, _    := g.Cfg().Get(ctx, prefix+".region")
	accessKey, _ := g.Cfg().Get(ctx, prefix+".accessKey")
	secretKey, _ := g.Cfg().Get(ctx, prefix+".secretKey")
	baseURL, _   := g.Cfg().Get(ctx, prefix+".baseUrl")
	useSSLVal, _ := g.Cfg().Get(ctx, prefix+".useSSL")
	pathStyleVal, _ := g.Cfg().Get(ctx, prefix+".pathStyle")

	useSSL := true
	if !useSSLVal.IsNil() {
		useSSL = useSSLVal.Bool()
	}

	opts := &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey.String(), secretKey.String(), ""),
		Secure: useSSL,
		Region: region.String(),
	}
	if pathStyleVal.Bool() {
		opts.BucketLookup = minio.BucketLookupPath
	}

	client, err := minio.New(endpoint.String(), opts)
	if err != nil {
		g.Log().Warningf(ctx, "[storage/s3] init client failed: %v", err)
	}

	return &s3Uploader{
		client:  client,
		bucket:  bucket.String(),
		baseURL: strings.TrimRight(baseURL.String(), "/"),
	}
}

func (u *s3Uploader) Upload(ctx context.Context, file *ghttp.UploadFile, opts UploadOptions) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	subDir  := renderPathTemplate(opts.PathTemplate, opts.Category)
	ext     := strings.ToLower(filepath.Ext(file.Filename))
	objKey  := subDir + "/" + randomHex(16) + ext

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	info, err := u.client.PutObject(ctx, u.bucket, objKey, src, file.Size,
		minio.PutObjectOptions{ContentType: mimeType})
	if err != nil {
		return nil, fmt.Errorf("s3 put object: %w", err)
	}

	cdnURL := u.baseURL + "/" + objKey

	return &UploadResult{
		StorageType: 2,
		StorageKey:  objKey,
		CdnUrl:      cdnURL,
		MimeType:    mimeType,
		FileSize:    info.Size,
	}, nil
}

func (u *s3Uploader) SaveBytes(ctx context.Context, data []byte, ext, prefix, originalKey string) (string, error) {
	dir  := filepath.ToSlash(filepath.Dir(originalKey))
	base := strings.TrimSuffix(filepath.Base(originalKey), filepath.Ext(originalKey))
	objKey := dir + "/" + base + "_" + prefix + ext

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	_, err := u.client.PutObject(ctx, u.bucket, objKey, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: mimeType})
	if err != nil {
		return "", fmt.Errorf("s3 save thumbnail: %w", err)
	}
	return u.baseURL + "/" + objKey, nil
}

func (u *s3Uploader) Delete(ctx context.Context, storageKey string) error {
	return u.client.RemoveObject(ctx, u.bucket, storageKey, minio.RemoveObjectOptions{})
}
