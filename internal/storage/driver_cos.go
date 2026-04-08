package storage

// Tencent Cloud COS storage driver.
// Config keys under storage.backends.<name>:
//
//	type:      "cos"
//	bucket:    "my-bucket-1234567890"
//	region:    "ap-beijing"
//	secretId:  "xxx"
//	secretKey: "xxx"
//	baseUrl:   "https://my-bucket-1234567890.cos.ap-beijing.myqcloud.com"   # optional

import (
	"bytes"
	"context"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	cos "github.com/tencentyun/cos-go-sdk-v5"
)

type cosUploader struct {
	client  *cos.Client
	baseURL string
}

func newCOSUploaderAt(ctx context.Context, prefix string) *cosUploader {
	bucket, _    := g.Cfg().Get(ctx, prefix+".bucket")
	region, _    := g.Cfg().Get(ctx, prefix+".region")
	secretID, _  := g.Cfg().Get(ctx, prefix+".secretId")
	secretKey, _ := g.Cfg().Get(ctx, prefix+".secretKey")
	baseURLCfg, _ := g.Cfg().Get(ctx, prefix+".baseUrl")

	bucketURL := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucket.String(), region.String())
	u, _ := url.Parse(bucketURL)

	client := cos.NewClient(&cos.BaseURL{BucketURL: u}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID.String(),
			SecretKey: secretKey.String(),
		},
	})

	baseURL := strings.TrimRight(baseURLCfg.String(), "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(bucketURL, "/")
	}

	return &cosUploader{client: client, baseURL: baseURL}
}

func (u *cosUploader) Upload(ctx context.Context, file *ghttp.UploadFile, opts UploadOptions) (*UploadResult, error) {
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

	_, err = u.client.Object.Put(ctx, objKey, src, &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   mimeType,
			ContentLength: file.Size,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("cos put object: %w", err)
	}

	return &UploadResult{
		StorageType: 4,
		StorageKey:  objKey,
		CdnUrl:      u.baseURL + "/" + objKey,
		MimeType:    mimeType,
		FileSize:    file.Size,
	}, nil
}

func (u *cosUploader) SaveBytes(ctx context.Context, data []byte, ext, prefix, originalKey string) (string, error) {
	dir    := filepath.ToSlash(filepath.Dir(originalKey))
	base   := strings.TrimSuffix(filepath.Base(originalKey), filepath.Ext(originalKey))
	objKey := dir + "/" + base + "_" + prefix + ext

	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "image/jpeg"
	}

	_, err := u.client.Object.Put(ctx, objKey, bytes.NewReader(data), &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentType:   mimeType,
			ContentLength: int64(len(data)),
		},
	})
	if err != nil {
		return "", fmt.Errorf("cos save thumbnail: %w", err)
	}
	return u.baseURL + "/" + objKey, nil
}

func (u *cosUploader) Delete(ctx context.Context, storageKey string) error {
	_, err := u.client.Object.Delete(ctx, storageKey)
	return err
}
