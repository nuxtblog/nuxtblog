package storage

import (
	"context"
	"io"

	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/sdk"
)

// sdkAdapterBridge wraps a sdk.StorageAdapter to implement storage.Uploader,
// bridging the plugin SDK interface to the internal storage interface.
type sdkAdapterBridge struct {
	adapter     sdk.StorageAdapter
	storageType int
}

func (b *sdkAdapterBridge) Upload(ctx context.Context, file *ghttp.UploadFile, opts UploadOptions) (*UploadResult, error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}

	uf := &sdk.UploadFile{
		Filename: file.Filename,
		Size:     file.Size,
		Reader:   newBytesReadSeeker(data),
	}
	sdkOpts := sdk.UploadOpts{
		Category:     opts.Category,
		PathTemplate: opts.PathTemplate,
	}
	res, err := b.adapter.Upload(ctx, uf, sdkOpts)
	if err != nil {
		return nil, err
	}
	return &UploadResult{
		StorageType: res.StorageType,
		StorageKey:  res.StorageKey,
		CdnUrl:      res.CdnUrl,
		MimeType:    res.MimeType,
		FileSize:    res.FileSize,
		Width:       res.Width,
		Height:      res.Height,
		Duration:    res.Duration,
		Variants:    res.Variants,
	}, nil
}

func (b *sdkAdapterBridge) Delete(ctx context.Context, storageKey string) error {
	return b.adapter.Delete(ctx, storageKey)
}

// SaveBytes delegates to sdk.ThumbAdapter if the adapter implements it.
func (b *sdkAdapterBridge) SaveBytes(ctx context.Context, data []byte, ext, prefix, originalKey string) (string, error) {
	if ta, ok := b.adapter.(sdk.ThumbAdapter); ok {
		return ta.SaveBytes(ctx, data, ext, prefix, originalKey)
	}
	return "", nil
}

// bytesReadSeeker wraps a byte slice as io.ReadSeeker.
type bytesReadSeeker struct {
	data   []byte
	offset int64
}

func newBytesReadSeeker(data []byte) *bytesReadSeeker {
	return &bytesReadSeeker{data: data}
}

func (b *bytesReadSeeker) Read(p []byte) (int, error) {
	if b.offset >= int64(len(b.data)) {
		return 0, io.EOF
	}
	n := copy(p, b.data[b.offset:])
	b.offset += int64(n)
	return n, nil
}

func (b *bytesReadSeeker) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		b.offset = offset
	case io.SeekCurrent:
		b.offset += offset
	case io.SeekEnd:
		b.offset = int64(len(b.data)) + offset
	}
	return b.offset, nil
}
