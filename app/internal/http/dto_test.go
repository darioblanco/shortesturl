package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestURLPayload_Bind(t *testing.T) {
	data := &URLPayload{URL: "http://valid.com"}
	err := data.Bind(&http.Request{})
	assert.NoError(t, err)
}

func TestURLPayload_Bind_InvalidURL(t *testing.T) {
	t.Parallel()
	tests := []struct {
		url string
	}{
		{"invalid"},
		{""},
		{"://noscheme.com"},
		{"https://"},
		{"weird://weirdscheme.com"},
	}
	for _, tt := range tests {
		tt := tt // NOTE: https://github.com/golang/go/wiki/CommonMistakes#using-goroutines-on-loop-iterator-variables
		t.Run(tt.url, func(t *testing.T) {
			t.Parallel()
			data := &URLPayload{URL: tt.url}
			assert.Error(t, data.Bind(&http.Request{}))
		})
	}
}

func TestURLPayload_Render(t *testing.T) {
	data := &URLPayload{URL: "http://valid.com"}
	rr := httptest.NewRecorder()
	err := data.Render(rr, &http.Request{})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rr.Code)
}
