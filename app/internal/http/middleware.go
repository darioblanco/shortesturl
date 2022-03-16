package http

import (
	"context"
	"net/http"
	"time"

	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/go-chi/chi/v5/middleware"
)

// ContextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
// See https://github.com/go-chi/chi/blob/master/middleware/middleware.go
type ContextKey struct {
	Name string
}

// loggerCtxKey is the key that holds the logger information in a request context.
var loggerCtxKey = &ContextKey{Name: "Logger"}

// LoggerMW middleware is used to call the injected logger in each request
func LoggerMW(logger logging.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			t1 := time.Now()
			defer func() {
				logger.Info("Served",
					"ip", r.RemoteAddr,
					"latency", time.Since(t1),
					"path", r.URL.Path,
					"protocol", r.Proto,
					"referer", r.Header.Get("Referer"),
					"requestId", middleware.GetReqID(r.Context()),
					"size", ww.BytesWritten(),
					"status", ww.Status(),
					"userAgent", r.Header.Get("User-Agent"),
				)
			}()
			next.ServeHTTP(ww, WithLoggerMW(r, logger))
		}
		return http.HandlerFunc(fn)
	}
}

// WithAdminMW sets the in-context logger for a request.
func WithLoggerMW(r *http.Request, logger logging.Logger) *http.Request {
	r = r.WithContext(context.WithValue(r.Context(), loggerCtxKey, logger))
	return r
}
