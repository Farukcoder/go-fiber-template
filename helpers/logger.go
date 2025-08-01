package helpers

import (
	"fmt"
	"garma_track/models"
	"garma_track/requests"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

// AsyncLogger handles database logging
type AsyncLogger struct {
	db      *gorm.DB
	channel chan requests.LogEntry
}

// NewAsyncLogger creates a new async logger instance
func NewAsyncLogger(db *gorm.DB) *AsyncLogger {
	return &AsyncLogger{
		db:      db,
		channel: make(chan requests.LogEntry, 100), // Buffered channel
	}
}

// ProcessLog processes log entries from the channel
func (logger *AsyncLogger) ProcessLog() {
	Info("🚀 Starting asynchronous logger...")

	for logEntry := range logger.channel {
		Debug(fmt.Sprintf("Processing log entry: %s %s", logEntry.Method, logEntry.URL))

		dbLog := models.Log{
			Method:          logEntry.Method,
			URL:             logEntry.URL,
			RequestBody:     logEntry.RequestBody,
			ResponseBody:    logEntry.ResponseBody,
			RequestHeaders:  logEntry.RequestHeaders,
			ResponseHeaders: logEntry.ResponseHeaders,
			StatusCode:      logEntry.StatusCode,
			CreatedAt:       logEntry.CreatedAt,
		}

		if err := logger.db.Create(&dbLog).Error; err != nil {
			Error(fmt.Sprintf("Failed to insert log entry: %v", err), nil)
		} else {
			Debug(fmt.Sprintf("Inserted log entry: %s %s", dbLog.Method, dbLog.URL))
		}
	}
}

// Log pushes a log entry into the channel
func (logger *AsyncLogger) Log(entry requests.LogEntry) {
	logger.channel <- entry
}

var (
	logFile     *os.File
	currentDate string
	logMutex    sync.Mutex
)

// LogLevel represents different logging levels
type LogLevel string

const (
	INFO    LogLevel = "✨ INFO"
	WARNING LogLevel = "⚠️ WARNING"
	ERROR   LogLevel = "❌ ERROR"
	DEBUG   LogLevel = "🐛 DEBUG"
)

// InitLogger initializes the logger
func InitLogger() error {
	logMutex.Lock()
	defer logMutex.Unlock()

	if err := createLogsDirectory(); err != nil {
		return err
	}

	return rotateLogFile()
}

// createLogsDirectory creates the logs directory if it doesn't exist
func createLogsDirectory() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	logsDir := filepath.Join(dir, "logs")
	if _, err := os.Stat(logsDir); os.IsNotExist(err) {
		return os.MkdirAll(logsDir, 0755)
	}
	return nil
}

// rotateLogFile creates or rotates to a new log file for the current date
func rotateLogFile() error {
	today := time.Now().Format("2006-01-02")
	if today == currentDate && logFile != nil {
		return nil
	}

	// Close existing log file if it exists
	if logFile != nil {
		logFile.Close()
	}

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	logPath := filepath.Join(dir, "logs", fmt.Sprintf("%s.log", today))
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	logFile = file
	currentDate = today

	// Create a multi-writer that writes to both file and console
	multiWriter := io.MultiWriter(file, os.Stdout)
	log.SetOutput(multiWriter)
	return nil
}

// Info logs an info level message
func Info(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Info("✨ " + message)
}

// Warning logs a warning level message
func Warning(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	log.Warn("⚠️ " + message)
}

// Error logs an error level message
func Error(format string, err error) {
	if err != nil {
		log.Error("❌ " + format + ": " + err.Error())
	} else {
		log.Error("❌ " + format)
	}
}

// Debug logs a debug level message
func Debug(format string, v ...interface{}) {
	if os.Getenv("GO_ENV") != "production" {
		message := fmt.Sprintf(format, v...)
		log.Debug("🐛 " + message)
	}
}
