package sdk

// Request is the HTTP request object passed to plugin route handlers.
// Extracted from *ghttp.Request to keep the SDK free of framework deps.
type Request struct {
	Method   string            `json:"method"`
	Path     string            `json:"path"`
	Query    map[string]string `json:"query"`
	Body     any               `json:"body"`
	Headers  map[string]string `json:"headers"`
	UserID   int               `json:"userId,omitempty"`
	UserRole int               `json:"userRole,omitempty"`
}

// Response is the HTTP response returned by plugin route handlers.
type Response struct {
	Status  int               `json:"status"`
	Body    any               `json:"body"`
	Headers map[string]string `json:"headers,omitempty"`
}
