package http

import (
	"errors"
	"net/http"
	"net/url"

	"github.com/go-chi/render"
)

// A URLPayload defines the JSON payload for sending and receiving urls
type URLPayload struct {
	URL       string  `json:"url"`
	ParsedURL url.URL `json:"-"`
}

// Bind validates the incoming request payload
func (ur *URLPayload) Bind(r *http.Request) error {
	u, err := url.Parse(ur.URL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		// I have considered that the url shortener should only accept HTTP or HTTPS
		// schemes, because I think that the next step would be to perform an
		// automatic redirect of these endpoints
		return errors.New("invalid http/https url format")
	}
	ur.ParsedURL = *u
	return nil
}

// Render defines the HTTP status code to 200
func (ur *URLPayload) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusOK)
	return nil
}

// A LongURL struct for the Swagger documentation
type LongURL struct {
	URL string `json:"url" example:"https://github.com/darioblanco"`
}

// A ShortURL struct for the Swagger documentation
type ShortURL struct {
	URL string `json:"url" example:"http://localhost:3000/64fc5e"`
}

// A BadRequest error struct for the Swagger documentation
type BadRequest struct {
	StatusText string `json:"status" example:"Bad Request"`
	ErrorText  string `json:"error,omitempty" example:"invalid http/https url format"`
}

// A NotFound error struct for the Swagger documentation
type NotFound struct {
	StatusText string `json:"status" example:"Not Found"`
	ErrorText  string `json:"error,omitempty" example:"long url not found"`
}

// A InternalServerError error struct for the Swagger documentation
type InternalServerError struct {
	StatusText string `json:"status" example:"Internal Server Error"`
	ErrorText  string `json:"error,omitempty" example:"oops, something went wrong in our side"`
}
