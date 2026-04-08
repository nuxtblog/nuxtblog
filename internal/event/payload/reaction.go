package payload

// ReactionAdded is delivered when a user adds a like or bookmark.
type ReactionAdded struct {
	UserID     int64
	ObjectType string // "post"
	ObjectID   int64
	Type       string // "like" | "bookmark"
}

// ReactionRemoved is delivered when a user removes a like or bookmark.
type ReactionRemoved struct {
	UserID     int64
	ObjectType string
	ObjectID   int64
	Type       string // "like" | "bookmark"
}
