package consts

// Storage driver type constants — used in config (storage.backends.<name>.type)
// and in the backend factory switch in storage/uploader.go.
// Cloud drivers (S3, OSS, COS) are now provided by builtin plugins.
const (
	StorageDriverLocal = "local" // local filesystem (default)
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

// BuiltinStorageDrivers lists driver types always available in the core.
// Plugin-provided drivers (S3, OSS, COS) register dynamically via the adapter registry.
var BuiltinStorageDrivers = []StorageDriverDef{
	{Type: StorageDriverLocal, LabelZh: "本地存储", LabelEn: "Local Filesystem"},
}
