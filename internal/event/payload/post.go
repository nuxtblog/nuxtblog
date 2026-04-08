// Package payload contains typed payload structs for each event.
// Import this package alongside the event package to perform type assertions:
//
//	event.OnAsync(event.PostPublished, func(ctx context.Context, e event.Event) error {
//	    p := e.Payload.(payload.PostPublished)
//	    ...
//	})
package payload

// PostCreated is delivered when any post is first inserted (any status).
type PostCreated struct {
	PostID   int64
	AuthorID int64
	Title    string
	Slug     string
	Excerpt  string
	// PostType: 0=post 1=page
	PostType int
	// Status: 0=draft 1=published 2=trash
	Status int
}

// PostUpdated is delivered when a post's fields change.
type PostUpdated struct {
	PostID   int64
	AuthorID int64
	Title    string
	Slug     string
	Excerpt  string
	PostType int
	Status   int
}

// PostPublished is delivered when a post transitions to the "published" status.
type PostPublished struct {
	PostID   int64
	AuthorID int64
	Title    string
	Slug     string
	Excerpt  string
	PostType int
}

// PostDeleted is delivered when a post is soft-deleted or moved to trash.
type PostDeleted struct {
	PostID   int64
	AuthorID int64
	Title    string
	Slug     string
	PostType int
}

// PostViewed is delivered each time a post's view counter is incremented.
// UserID is 0 for anonymous visitors.
type PostViewed struct {
	PostID int64
	UserID int64
}
