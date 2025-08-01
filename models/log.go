package models

import (
	"time"

	"gorm.io/gorm"
)

// Log represents a database log entry
type Log struct {
	gorm.Model
	Method          string    `json:"method"`
	URL             string    `json:"url"`
	RequestBody     string    `json:"request_body"`
	ResponseBody    string    `json:"response_body"`
	RequestHeaders  string    `json:"request_headers"`
	ResponseHeaders string    `json:"response_headers"`
	StatusCode      int       `json:"status_code"`
	CreatedAt       time.Time `json:"created_at"`
}
