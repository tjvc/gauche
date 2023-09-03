package middleware

import (
	"net/http"
	"time"

	"golang.org/x/exp/slog"
)

type responseWriterWithStatus struct {
	http.ResponseWriter
	Status int
}

func (w *responseWriterWithStatus) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

func Log(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writerWithStatus := &responseWriterWithStatus{ResponseWriter: w}
		next.ServeHTTP(writerWithStatus, r)
		end := time.Now()

		var level slog.Level
		var msg string

		if writerWithStatus.Status == http.StatusInternalServerError {
			level = slog.LevelError
			msg = "Request failed"
		} else {
			level = slog.LevelInfo
			msg = "Request processed"
		}

		logger.Log(
			nil,
			level,
			msg,
			"method", r.Method,
			"path", r.URL.Path,
			"status", writerWithStatus.Status,
			"latency", end.Sub(start),
		)
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
