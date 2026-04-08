package payload

// OptionUpdated is delivered when an admin writes a site option via the API.
type OptionUpdated struct {
	Key   string
	Value string
}
