package http

import (
	"context"
	"net/http"

	"github.com/darioblanco/shortesturl/app/internal/cache"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

// NewRouter creates a chi router conforming the Handler interface
func NewRouter(
	ctx context.Context, conf *config.Values, logger logging.Logger, cache cache.Cache,
) (http.Handler, error) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Compress(5))
	r.Use(LoggerMW(logger))
	r.Use(middleware.URLFormat)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Use(middleware.Heartbeat("/health"))

	r.Mount("/docs", docs{config: conf}.Router())
	r.Mount("/", api{cache: cache, config: conf, logger: logger}.Router())

	return r, nil
}
