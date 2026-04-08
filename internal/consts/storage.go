package consts

// Storage driver type constants — used in config (storage.backends.<name>.type)
// and in the backend factory switch in storage/uploader.go.
// To add a new driver: implement storage.Uploader, add a case in buildBackend(), then add an entry below.
const (
	StorageDriverLocal = "local" // local filesystem (default)
	StorageDriverS3    = "s3"   // S3-compatible (AWS, MinIO, Cloudflare R2 …)
	StorageDriverOSS   = "oss"  // Aliyun OSS
	StorageDriverCOS   = "cos"  // Tencent COS
)

// Upload path template placeholders — used in media_upload_path option.
// Combine freely, e.g. "{year}/{month}" or "{year}/{month}/{day}/{category}".
const (
	StoragePathTemplateYearMonth    = "{year}/{month}"
	StoragePathTemplateYearMonthDay = "{year}/{month}/{day}"
	StoragePathTemplateFlat         = "files"

	StoragePathPlaceholderYear     = "{year}"
	StoragePathPlaceholderMonth    = "{month}"
	StoragePathPlaceholderDay      = "{day}"
	StoragePathPlaceholderCategory = "{category}"
)

// StorageDriverDef describes a storage driver type.
type StorageDriverDef struct {
	Type    string // matches config type value
	LabelZh string
	LabelEn string
}

// BuiltinStorageDrivers lists all driver types implemented in this codebase.
// The admin settings page uses this to validate and label storage backend types.
var BuiltinStorageDrivers = []StorageDriverDef{
	{Type: StorageDriverLocal, LabelZh: "本地存储",    LabelEn: "Local Filesystem"},
	{Type: StorageDriverS3,    LabelZh: "S3 / MinIO", LabelEn: "S3 / MinIO"},
	{Type: StorageDriverOSS,   LabelZh: "阿里云 OSS", LabelEn: "Aliyun OSS"},
	{Type: StorageDriverCOS,   LabelZh: "腾讯云 COS", LabelEn: "Tencent COS"},
}
