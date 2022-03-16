package http

import (
	"context"
	"testing"

	"github.com/darioblanco/shortesturl/app/internal/cache"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewRouter(t *testing.T) {
	r, err := NewRouter(
		context.Background(),
		&config.Values{
			HttpScheme: "http",
			HttpHost:   "localhost",
			HttpPort:   3000,
		},
		logging.NewTest(t),
		cache.NewTest(),
	)
	assert.NoError(t, err)
	assert.Len(t, r.(*chi.Mux).Middlewares(), 8)
	assert.Len(t, r.(*chi.Mux).Routes(), 2)
}
