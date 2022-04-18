package logging

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogEntry struct {
	Timestamp time.Time
	Status    int
	Method    string
	Path      string
	Latency   time.Duration
}

type Logger interface {
	Write(LogEntry)
}

type jsonLogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"status"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Latency    int64     `json:"latency_μs"`
}

type JSONLogger struct{}

func (logger JSONLogger) Write(logEntry LogEntry) {
	jsonLogEntry := jsonLogEntry{
		Timestamp:  logEntry.Timestamp,
		StatusCode: logEntry.Status,
		Method:     logEntry.Method,
		Path:       logEntry.Path,
		Latency:    logEntry.Latency.Microseconds(),
	}

	marshalledJSONLogEntry, _ := json.Marshal(jsonLogEntry)
	fmt.Println(string(marshalledJSONLogEntry))
}
