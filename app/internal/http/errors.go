package http

import (
	"net/http"

	"github.com/go-chi/render"
)

// ErrHTTPResponse renderer type for handling all sorts of errors.
//
// In the best case scenario, the excellent github.com/pkg/errors package
// helps reveal information on the error, setting it on Err, and in the Render()
// method, using it to set the application-specific error code in AppCode.
type ErrHTTPResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	ErrorText  string `json:"error,omitempty"` // application-level error message
}

// Render defines the HTTP status code based on its inherent error
func (e *ErrHTTPResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

// ErrBadRequest returns a 400 and the error message
// The message from the error passed as parameter IS shown to the end user
func ErrBadRequest(err error) render.Renderer {
	return &ErrHTTPResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     http.StatusText(http.StatusBadRequest),
		ErrorText:      err.Error(),
	}
}

// ErrNotFound returns a 404 and the error message
// The message from the error passed as parameter IS shown to the end user
func ErrNotFound(err error) render.Renderer {
	return &ErrHTTPResponse{
		Err:            err,
		HTTPStatusCode: http.StatusNotFound,
		StatusText:     http.StatusText(http.StatusNotFound),
		ErrorText:      err.Error(),
	}
}

// ErrInternalServerError returns a 500 and a generic message
// The message from the error passed as parameter IS NOT shown to the end user
func ErrInternalServerError(err error) render.Renderer {
	return &ErrHTTPResponse{
		Err:            err,
		HTTPStatusCode: http.StatusInternalServerError,
		StatusText:     http.StatusText(http.StatusInternalServerError),
		ErrorText:      "oops, something went wrong in our side",
	}
}
