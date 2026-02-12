package logger

import (
	"fmt"
	"sync"
	"time"
)

// LogEntry represents a single log line for the frontend
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"` // INFO, WARN, ERROR, SUCCESS
	Message   string `json:"message"`
}

var (
	logs []LogEntry
	mu   sync.Mutex
)

// Init initializes the logger buffer
func Init() {
	logs = make([]LogEntry, 0)
}

// Internal helper to add log
func add(level, format string, args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()

	msg := fmt.Sprintf(format, args...)
	entry := LogEntry{
		Timestamp: time.Now().Format("15:04:05"),
		Level:     level,
		Message:   msg,
	}

	// Append and keep only last 50 lines
	logs = append(logs, entry)
	if len(logs) > 50 {
		logs = logs[1:]
	}

	// Mirror to Standard Output for debugging
	color := "\033[0m" // Reset
	switch level {
	case "ERROR":
		color = "\033[31m" // Red
	case "WARN":
		color = "\033[33m" // Yellow
	case "SUCCESS":
		color = "\033[32m" // Green
	case "INFO":
		color = "\033[34m" // Blue
	case "DEBUG":
		color = "\033[90m" // Grey
	}
	fmt.Printf("%s[%s] %s\033[0m\n", color, level, msg)
}

// Public Methods
func Info(format string, args ...interface{}) {
	add("INFO", format, args...)
}

func Debug(format string, args ...interface{}) {
	add("DEBUG", format, args...)
}

func Warn(format string, args ...interface{}) {
	add("WARN", format, args...)
}

func Error(format string, args ...interface{}) {
	add("ERROR", format, args...)
}

func Success(format string, args ...interface{}) {
	add("SUCCESS", format, args...)
}

func GetLogs() []LogEntry {
	mu.Lock()
	defer mu.Unlock()
	// Return a copy to prevent race conditions
	dst := make([]LogEntry, len(logs))
	copy(dst, logs)
	return dst
}
