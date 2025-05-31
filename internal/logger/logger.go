package logger

import (
	"fmt"
	"time"
)

// ANSI color codes
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

// Log levels
const (
	LevelInfo    = "INFO"
	LevelError   = "ERROR"
	LevelSuccess = "SUCCESS"
	LevelWarn    = "WARN"
)

// getColorByLevel returns the ANSI color code for a log level
func getColorByLevel(level string) string {
	switch level {
	case LevelInfo:
		return Blue
	case LevelError:
		return Red
	case LevelSuccess:
		return Green
	case LevelWarn:
		return Yellow
	default:
		return Reset
	}
}

// formatMessage formats a log message with timestamp and level
func formatMessage(level, msg string) string {
	timestamp := time.Now().Format("2006/01/02 15:04:05")
	color := getColorByLevel(level)
	return fmt.Sprintf("%s[%s] %s%s%s %s",
		Purple, timestamp, color, level, Reset, msg)
}

// Info logs an informational message
func Info(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(formatMessage(LevelInfo, msg))
}

// Error logs an error message
func Error(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(formatMessage(LevelError, msg))
}

// Success logs a success message
func Success(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(formatMessage(LevelSuccess, msg))
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Println(formatMessage(LevelWarn, msg))
}

// CrawlStatus logs the current crawling status
func CrawlStatus(url string, linksFound int, totalPages, queueSize int) {
	msg := fmt.Sprintf("Crawled: %s%s%s | Links found: %s%d%s | Total pages: %s%d%s | Queue size: %s%d%s",
		Cyan, url, Reset,
		Green, linksFound, Reset,
		Yellow, totalPages, Reset,
		Purple, queueSize, Reset)
	fmt.Println(formatMessage(LevelInfo, msg))
}

// StorageStatus logs MongoDB storage operations
func StorageStatus(url string, isUpdate bool) {
	action := "Stored"
	if isUpdate {
		action = "Updated"
	}
	msg := fmt.Sprintf("%s page: %s%s%s", action, Cyan, url, Reset)
	fmt.Println(formatMessage(LevelSuccess, msg))
}
