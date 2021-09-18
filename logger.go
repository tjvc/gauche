package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
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

	marshalledJsonLogEntry, _ := json.Marshal(jsonLogEntry)
	fmt.Println(string(marshalledJsonLogEntry))
}

func loggingMiddleware(logger logger) func(c *gin.Context) {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		end := time.Now()

		logEntry := logEntry{
			timestamp: end,
			status:    c.Writer.Status(),
			method:    c.Request.Method,
			path:      c.Request.URL.Path,
			latency:   end.Sub(start),
		}

		logger.write(logEntry)
	}
}
