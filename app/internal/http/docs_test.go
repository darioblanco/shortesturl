package http

import (
	"testing"

	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/stretchr/testify/assert"
)

func TestDocs(t *testing.T) {
	router := docs{config: &config.Values{
		HttpScheme: "http",
		HttpHost:   "localhost",
		HttpPort:   3000,
	}}.Router()
	assert.Len(t, router.Routes(), 1)
}
