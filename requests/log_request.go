package requests

import "time"

// LogEntry represents a log entry in the system
type LogEntry struct {
	Method          string
	URL             string
	RequestBody     string
	ResponseBody    string
	RequestHeaders  string
	ResponseHeaders string
	StatusCode      int
	CreatedAt       time.Time
}
