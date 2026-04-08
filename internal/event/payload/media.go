package payload

// MediaUploaded is delivered after a file is successfully stored and its DB record created.
type MediaUploaded struct {
	MediaID    int64
	UploaderID int64
	Filename   string
	MimeType   string
	// FileSize in bytes
	FileSize int64
	// URL is the public CDN / storage URL of the original file.
	URL      string
	Category string
	// Width and Height are non-zero for images.
	Width  int
	Height int
}

// MediaDeleted is delivered after a media record and its storage file are removed.
type MediaDeleted struct {
	MediaID    int64
	UploaderID int64
	Filename   string
	MimeType   string
	Category   string
}
