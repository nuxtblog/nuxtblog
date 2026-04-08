package payload

// UserRegistered is delivered when a new user account is created.
type UserRegistered struct {
	UserID      int64
	Username    string
	Email       string
	DisplayName string
	Locale      string
	// Role: 0=subscriber 1=contributor 2=editor 3=admin
	Role int
}

// UserUpdated is delivered when a user's profile or account fields change.
type UserUpdated struct {
	UserID      int64
	Username    string
	Email       string
	DisplayName string
	Locale      string
	Role        int
	// Status: 0=active 1=inactive/banned
	Status int
}

// UserDeleted is delivered when a user account is soft-deleted.
type UserDeleted struct {
	UserID   int64
	Username string
	Email    string
}

// UserLoggedIn is delivered after a successful password or OAuth login.
type UserLoggedIn struct {
	UserID   int64
	Username string
	Email    string
	// Role: 0=subscriber 1=contributor 2=editor 3=admin
	Role int
}

// UserLoggedOut is delivered after a session is revoked (logout).
// UserID is 0 when the logout was by refresh-token only (user not resolved).
type UserLoggedOut struct {
	UserID int64
}
