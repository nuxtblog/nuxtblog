package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gtime"
)

// ----------------------------------------------------------------
//  Enums
// ----------------------------------------------------------------

type StorageType int

const (
	StorageTypeLocal    StorageType = 1
	StorageTypeS3       StorageType = 2
	StorageTypeOSS      StorageType = 3
	StorageTypeCOS      StorageType = 4
	StorageTypeExternal StorageType = 5 // external URL, no file stored
)

// ----------------------------------------------------------------
//  Shared output
// ----------------------------------------------------------------

type MediaItem struct {
	Id          int64       `json:"id"`
	UploaderId  int64       `json:"uploader_id"`
	StorageType StorageType `json:"storage_type"`
	CdnUrl      string      `json:"cdn_url"`
	Filename    string      `json:"filename"`
	MimeType    string      `json:"mime_type"`
	FileSize    int64       `json:"file_size"`
	Width       *int        `json:"width"`
	Height      *int        `json:"height"`
	Duration    *int        `json:"duration"`
	AltText     string      `json:"alt_text"`
	Title       string      `json:"title"`
	Category    string      `json:"category"`
	Variants    string      `json:"variants"` // JSON
	CreatedAt   *gtime.Time `json:"created_at"`
}

// ----------------------------------------------------------------
//  Upload
// ----------------------------------------------------------------

type MediaUploadReq struct {
	g.Meta  `path:"/medias/upload" method:"post" tags:"Media" summary:"Upload a file" mime:"multipart/form-data"`
	File     *ghttp.UploadFile `v:"required" dc:"file to upload"`
	AltText  string            `v:"max-length:255" dc:"alt text for image"`
	Title    string            `v:"max-length:255" dc:"media title"`
	Category string            `v:"max-length:32" dc:"category: avatar|cover|post|doc|moment|banner"`
}
type MediaUploadRes struct {
	*MediaItem `dc:"uploaded media"`
}

// ----------------------------------------------------------------
//  Delete
// ----------------------------------------------------------------

type MediaDeleteReq struct {
	g.Meta `path:"/medias/{id}" method:"delete" tags:"Media" summary:"Delete media"`
	Id     int64 `v:"required|min:1" dc:"media id"`
}
type MediaDeleteRes struct{}

// ----------------------------------------------------------------
//  Update
// ----------------------------------------------------------------

type MediaUpdateReq struct {
	g.Meta  `path:"/medias/{id}" method:"put" tags:"Media" summary:"Update media metadata"`
	Id      int64   `v:"required|min:1"  dc:"media id"`
	AltText *string `v:"max-length:255"  dc:"alt text"`
	Title   *string `v:"max-length:255"  dc:"title"`
}
type MediaUpdateRes struct{}

// ----------------------------------------------------------------
//  Get one
// ----------------------------------------------------------------

type MediaGetOneReq struct {
	g.Meta `path:"/medias/{id}" method:"get" tags:"Media" summary:"Get media by id"`
	Id     int64 `v:"required|min:1" dc:"media id"`
}
type MediaGetOneRes struct {
	*MediaItem `dc:"media"`
}

// ----------------------------------------------------------------
//  Get list
// ----------------------------------------------------------------

type MediaGetListReq struct {
	g.Meta      `path:"/medias" method:"get" tags:"Media" summary:"Get media list"`
	MimeType    *string `dc:"filter by mime type, e.g. image/jpeg"`
	UploaderId  *int64  `v:"min:1"         dc:"filter by uploader"`
	Category    *string `v:"max-length:32" dc:"filter by category: avatar|cover|post|doc|moment|banner"`
	StorageType *int    `dc:"filter by storage type: 1=local 2=S3 3=OSS 4=COS 5=external"`
	Page        int     `v:"min:1"         dc:"page number" d:"1"`
	Size        int     `v:"between:1,100" dc:"page size"   d:"20"`
}
type MediaGetListRes struct {
	List  []*MediaItem `json:"list"`
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Size  int          `json:"size"`
}

// ----------------------------------------------------------------
//  Link (external URL)
// ----------------------------------------------------------------

// MediaLinkReq registers an external URL as a media record without uploading a file.
type MediaLinkReq struct {
	g.Meta   `path:"/medias/link" method:"post" tags:"Media" summary:"Register external media URL"`
	Url      string `v:"required|url|max-length:2048" dc:"external image or file URL"`
	Title    string `v:"max-length:255"               dc:"media title"`
	AltText  string `p:"alt_text" v:"max-length:255"  dc:"alt text for image"`
	Category string `v:"max-length:32"                dc:"category: avatar|cover|post|doc|moment|banner"`
}
type MediaLinkRes struct {
	*MediaItem `dc:"created media record"`
}

// ----------------------------------------------------------------
//  Stats
// ----------------------------------------------------------------

type MediaGetStatsReq struct {
	g.Meta `path:"/medias/stats" method:"get" tags:"Media" summary:"Get media stats by type"`
}
type MediaGetStatsRes struct {
	Total int `json:"total"`
	Image int `json:"image"`
	Video int `json:"video"`
	Audio int `json:"audio"`
	Other int `json:"other"`
}

// ----------------------------------------------------------------
//  Localize (download external URL → local storage)
// ----------------------------------------------------------------

type MediaLocalizeReq struct {
	g.Meta `path:"/medias/{id}/localize" method:"post" tags:"Media" summary:"Download external media URL to local storage"`
	Id     int64 `v:"required|min:1" dc:"media id"`
}
type MediaLocalizeRes struct {
	*MediaItem `dc:"updated media record"`
}
