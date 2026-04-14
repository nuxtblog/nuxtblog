package media

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"

	"github.com/disintegration/imaging"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"github.com/nuxtblog/nuxtblog/internal/consts"
	"github.com/nuxtblog/nuxtblog/internal/storage"
)

// getThumbSizes reads thumbnail/cover/content size presets from options,
// falling back to built-in defaults.
func getThumbSizes(ctx context.Context) (thumbnail, cover, content consts.ThumbSize) {
	thumbnail = consts.DefaultThumbThumbnail
	cover = consts.DefaultThumbCover
	content = consts.DefaultThumbContent
	_ = getOptionJSON(ctx, "media_thumbnail", &thumbnail)
	_ = getOptionJSON(ctx, "media_cover_thumb", &cover)
	_ = getOptionJSON(ctx, "media_content_thumb", &content)
	return
}

// generateThumbsFromFile reads image bytes from the upload file, resizes/crops them,
// and saves thumbnails via the ThumbSaver interface.
// Errors are non-fatal: logged and skipped.
func generateThumbsFromFile(ctx context.Context, up storage.Uploader, file *ghttp.UploadFile, originalKey string, thumbnail, cover, content consts.ThumbSize) map[string]string {
	ts, ok := up.(storage.ThumbSaver)
	if !ok {
		return nil
	}

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

	return generateThumbVariants(ctx, ts, img, originalKey, thumbnail, cover, content)
}

// generateThumbsFromBytes decodes image data from raw bytes, resizes/crops them,
// and saves thumbnails via the ThumbSaver interface.
// Returns nil if the uploader doesn't support ThumbSaver or the image can't be decoded.
func generateThumbsFromBytes(ctx context.Context, up storage.Uploader, data []byte, originalKey string, thumbnail, cover, content consts.ThumbSize) map[string]string {
	ts, ok := up.(storage.ThumbSaver)
	if !ok {
		return nil
	}

	img, err := imaging.Decode(bytes.NewReader(data))
	if err != nil {
		g.Log().Warningf(ctx, "thumbnail: decode image: %v", err)
		return nil
	}

	return generateThumbVariants(ctx, ts, img, originalKey, thumbnail, cover, content)
}

// generateThumbVariants creates thumbnail/cover/content variants from a decoded image.
func generateThumbVariants(ctx context.Context, ts storage.ThumbSaver, img image.Image, originalKey string, thumbnail, cover, content consts.ThumbSize) map[string]string {
	variants := make(map[string]string)

	saveVariant := func(name string, w, h int) {
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

	saveVariant("thumbnail", thumbnail.Width, thumbnail.Height)
	saveVariant("cover", cover.Width, cover.Height)
	saveVariant("content", content.Width, content.Height)

	return variants
}
