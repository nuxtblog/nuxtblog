package payload

// CommentCreated is delivered when a comment is successfully inserted.
type CommentCreated struct {
	CommentID int64
	// Status: 0=pending 1=approved 2=spam
	Status int

	ObjectType  string // "post"
	ObjectID    int64
	ObjectTitle string
	ObjectSlug  string

	// PostAuthorID is the author of the object being commented on.
	PostAuthorID int64

	// ParentID is non-nil for replies.
	ParentID *int64
	// ParentAuthorID is the author of the parent comment (non-zero for replies).
	ParentAuthorID int64

	AuthorID    int64 // 0 for anonymous comments
	AuthorName  string
	AuthorEmail string

	Content string
}

// CommentDeleted is delivered when a comment is soft-deleted.
type CommentDeleted struct {
	CommentID  int64
	ObjectType string
	ObjectID   int64
}

// CommentStatusChanged is delivered when an admin moderates a comment.
type CommentStatusChanged struct {
	CommentID  int64
	ObjectType string
	ObjectID   int64
	// OldStatus / NewStatus: 0=pending 1=approved 2=spam
	OldStatus int
	NewStatus int
	// ModeratorID is the admin who changed the status (0 if unknown).
	ModeratorID int64
}

// CommentApproved is a convenience event fired when a comment transitions
// to approved status (NewStatus == 2 in CommentStatusChanged).
type CommentApproved struct {
	CommentID   int64
	ObjectType  string
	ObjectID    int64
	ModeratorID int64
}
