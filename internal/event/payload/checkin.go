package payload

// CheckinDone is delivered when a user completes a check-in (including repeat check-ins).
type CheckinDone struct {
	UserID           int64
	Streak           int
	AlreadyCheckedIn bool
}
