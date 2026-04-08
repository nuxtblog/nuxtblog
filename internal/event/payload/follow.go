package payload

// UserFollowed is delivered when a user follows another user.
type UserFollowed struct {
	// FollowerID is the user who performed the follow action.
	FollowerID     int64
	FollowerName   string
	FollowerAvatar string

	// FollowingID is the user who was followed (the notification recipient).
	FollowingID int64
}
