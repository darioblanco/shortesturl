package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
)

func TestErrHTTPResponseRender(t *testing.T) {
	data := &ErrHTTPResponse{
		Err:            errors.New("Fake"),
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      "Fake error",
	}
	r, _ := http.NewRequest(http.MethodGet, "", nil)
	rr := httptest.NewRecorder()
	err := data.Render(rr, r)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, r.Context().Value(render.StatusCtxKey).(int))
}

func TestErrBadRequest(t *testing.T) {
	err := errors.New("Unknown error")
	res := ErrBadRequest(err)
	assert.Equal(t, &ErrHTTPResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
	}, res)
}

func TestErrNotFound(t *testing.T) {
	err := errors.New("Unknown error")
	res := ErrNotFound(err)
	assert.Equal(t, &ErrHTTPResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		ErrorText:      err.Error(),
	}, res)
}

func TestErrInternalServerError(t *testing.T) {
	err := errors.New("something really bad")
	res := ErrInternalServerError(err)
	assert.Equal(t, &ErrHTTPResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
		ErrorText:      "oops, something went wrong in our side",
	}, res)
}
