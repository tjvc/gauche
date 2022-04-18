package middleware

import (
	"net/http"
	"time"

	"github.com/tjvc/gauche/internal/logging"
)

type responseWriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func (w *responseWriterWithStatus) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func Log(logger logging.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writerWithStatus := &responseWriterWithStatus{ResponseWriter: w}
		next.ServeHTTP(writerWithStatus, r)
		end := time.Now()

		logEntry := logging.LogEntry{
			Timestamp: end,
			Status:    writerWithStatus.Status,
			Method:    r.Method,
			Path:      r.URL.Path,
			Latency:   end.Sub(start),
		}

		logger.Write(logEntry)
	})
}

func Recover(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(w, r)
	})
}
