package middleware

import (
	"net/http"
	"time"

	"github.com/alexmacdonald/simple-ipass/internal/logger"
)

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// RequestLogger logs HTTP requests with status codes, execution time, and metadata
// This provides observability for API performance and debugging
func RequestLogger(log *logger.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Wrap response writer to capture status code
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     http.StatusOK, // Default to 200
			}

			// Call next handler
			next.ServeHTTP(wrapped, r)

			// Calculate duration
			duration := time.Since(start)

			// Extract user ID if available (for authenticated requests)
			userID := ""
			if uid, ok := GetUserIDFromContext(r.Context()); ok {
				userID = uid
			}

			// Log the request with structured data
		level := logger.LevelInfo
		if wrapped.statusCode >= 500 {
			level = logger.LevelError
		} else if wrapped.statusCode >= 400 {
			level = logger.LevelWarn
		}

		logData := map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status_code": wrapped.statusCode,
			"duration_ms": duration.Milliseconds(),
			"duration":    duration.String(),
			"user_agent":  r.UserAgent(),
			"remote_addr": r.RemoteAddr,
			"user_id":     userID,
			"bytes_sent":  wrapped.written,
		}

		switch level {
		case logger.LevelError:
			log.Error("HTTP Request", logData)
		case logger.LevelWarn:
			log.Warn("HTTP Request", logData)
		default:
			log.Info("HTTP Request", logData)
		}
		})
	}
}

