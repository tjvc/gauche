package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type logEntry struct {
	timestamp time.Time
	status    int
	method    string
	path      string
	latency   time.Duration
}

type logger interface {
	write(logEntry)
}

type jsonLogEntry struct {
	Timestamp  time.Time `json:"timestamp"`
	StatusCode int       `json:"status"`
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	Latency    int64     `json:"latency_μs"`
}

type jsonLogger struct{}

func (logger jsonLogger) write(logEntry logEntry) {
	jsonLogEntry := jsonLogEntry{
		Timestamp:  logEntry.timestamp,
		StatusCode: logEntry.status,
		Method:     logEntry.method,
		Path:       logEntry.path,
		Latency:    logEntry.latency.Microseconds(),
	}

	marshalledJSONLogEntry, _ := json.Marshal(jsonLogEntry)
	fmt.Println(string(marshalledJSONLogEntry))
}

type responseWriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func (w *responseWriterWithStatus) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func loggingMiddleware(logger logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writerWithStatus := &responseWriterWithStatus{ResponseWriter: w}
		next.ServeHTTP(writerWithStatus, r)
		end := time.Now()

		logEntry := logEntry{
			timestamp: end,
			status:    writerWithStatus.Status,
			method:    r.Method,
			path:      r.URL.Path,
			latency:   end.Sub(start),
		}

		logger.write(logEntry)
	})
}
