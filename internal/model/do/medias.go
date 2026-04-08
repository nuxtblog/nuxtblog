// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// Medias is the golang structure of table medias for DAO operations like Where/Data.
type Medias struct {
	g.Meta      `orm:"table:medias, do:true"`
	Id          any         //
	UploaderId  any         //
	StorageType any         //
	StorageKey  any         //
	CdnUrl      any         //
	Filename    any         //
	MimeType    any         //
	FileSize    any         //
	Width       any         //
	Height      any         //
	Duration    any         //
	AltText     any         //
	Title       any         //
	Variants    any         //
	FileMeta    any         //
	CreatedAt   *gtime.Time //
	DeletedAt   *gtime.Time //
}
