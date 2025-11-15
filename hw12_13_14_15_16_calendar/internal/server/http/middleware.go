package internalhttp

import (
	"net/http"
	"time"

	"github.com/esivanov203/otus/hw12_13_14_15_calendar/internal/logger"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rw, r)

		s.logger.Info("HTTP request",
			logger.Fields{
				"ip":        r.RemoteAddr,
				"method":    r.Method,
				"path":      r.URL.Path,
				"proto":     r.Proto,
				"status":    rw.status,
				"latency":   time.Since(start).Milliseconds(),
				"userAgent": r.UserAgent(),
			})
	})
}
