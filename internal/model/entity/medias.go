// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// Medias is the golang structure for table medias.
type Medias struct {
	Id          int         `json:"id"          orm:"id"           description:""` //
	UploaderId  int         `json:"uploaderId"  orm:"uploader_id"  description:""` //
	StorageType int         `json:"storageType" orm:"storage_type" description:""` //
	StorageKey  string      `json:"storageKey"  orm:"storage_key"  description:""` //
	CdnUrl      string      `json:"cdnUrl"      orm:"cdn_url"      description:""` //
	Filename    string      `json:"filename"    orm:"filename"     description:""` //
	MimeType    string      `json:"mimeType"    orm:"mime_type"    description:""` //
	FileSize    int         `json:"fileSize"    orm:"file_size"    description:""` //
	Width       int         `json:"width"       orm:"width"        description:""` //
	Height      int         `json:"height"      orm:"height"       description:""` //
	Duration    int         `json:"duration"    orm:"duration"     description:""` //
	AltText     string      `json:"altText"     orm:"alt_text"     description:""` //
	Title       string      `json:"title"       orm:"title"        description:""` //
	Category    string      `json:"category"    orm:"category"     description:""` //
	Variants    string      `json:"variants"    orm:"variants"     description:""` //
	FileMeta    string      `json:"fileMeta"    orm:"file_meta"    description:""` //
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:""` //
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:""` //
}
