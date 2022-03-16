package http

import (
	"crypto/md5"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/darioblanco/shortesturl/app/internal/cache"
	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type api struct {
	cache  cache.Cache
	config *config.Values
	logger logging.Logger
}

func (rs api) Router() chi.Router {
	r := chi.NewRouter()

	r.Post("/encode", rs.Encode)
	r.Post("/decode", rs.Decode)

	return r
}

// Encode
// @Summary Encodes a URL to a shortened URL
// @Description Shorten a given URL, which can be decoded later using /decode
// @ID encode
// @Tags Shortener
// @Accept json
// @Produce json
// @Param url body LongURL true "The url to encode"
// @Success 200 {object} ShortURL "Long URL encoded successfully"
// @Failure 400 {object} BadRequest "Long URL has a wrong format"
// @Failure 500 {object} InternalServerError "Unexpected error in the backend"
// @Router /encode [post]
func (rs api) Encode(w http.ResponseWriter, r *http.Request) {
	data := &URLPayload{}
	if err := render.Bind(r, data); err != nil {
		rs.logger.Warn("long URL has a wrong format", "error", err)
		render.Render(w, r, ErrBadRequest(err))
		return
	}

	// Generate MD5 from url
	md5Url := fmt.Sprintf("%x", md5.Sum([]byte(data.URL)))
	rs.logger.Debug("MD5 encoded url",
		"url", data.URL,
		"MD5", md5Url,
	)

	// Shorten URL to the length defined in the application config
	success := false
	attempts := 0
	var err error
	var shortURLSlug string
	ctx := r.Context()
	for !success {
		shortURLSlug = md5Url[attempts : rs.config.UrlLength+attempts]
		// We could also send metrics to trigger a threshold if attempts get out of hand
		rs.logger.Debug("Attempting to store shortened url",
			"attempts", attempts,
			"shortUrlSlug", shortURLSlug,
		)
		// If success is false, it indicates a collision and the md5 has to be shifted
		success, err = rs.cache.SetIfNotExists(
			ctx,
			shortURLSlug,
			data.URL,
			time.Hour*rs.config.UrlExpirationInHours,
		)
		if err != nil {
			rs.logger.Error("unable to retrieve/store shortened url in cache", "error", err)
			render.Render(w, r, ErrInternalServerError(err))
			return
		}
		attempts += 1
	}
	var host string
	if (rs.config.HttpScheme == "https" && rs.config.HttpPort == 443) ||
		(rs.config.HttpScheme == "http" && rs.config.HttpPort == 80) {
		host = rs.config.HttpHost
	} else {
		host = fmt.Sprintf("%s:%d", rs.config.HttpHost, rs.config.HttpPort)
	}
	shortURL := url.URL{
		Scheme: rs.config.HttpScheme,
		Host:   host,
		Path:   string(shortURLSlug),
	}
	shortURLString := shortURL.String()
	rs.logger.Info("Encoded url",
		"longUrl", data.URL,
		"shortUrl", shortURLString,
	)
	render.Render(w, r, &URLPayload{URL: shortURLString})
}

// Decode
// @Summary Decodes a URL to a shortened URL
// @Description Revert a shortened URL to its original form
// @ID decode
// @Tags Shortener
// @Accept json
// @Produce json
// @Param url body ShortURL true "The url to decode"
// @Success 200 {object} LongURL "Short URL decoded successfully"
// @Failure 400 {object} BadRequest "Short URL has a wrong format"
// @Failure 404 {object} NotFound "Short URL has not a related long URL"
// @Failure 500 {object} InternalServerError "Unexpected error in the backend"
// @Router /decode [post]
func (rs api) Decode(w http.ResponseWriter, r *http.Request) {
	data := &URLPayload{}
	if err := render.Bind(r, data); err != nil {
		rs.logger.Warn("short URL has a wrong format", "error", err)
		render.Render(w, r, ErrBadRequest(err))
		return
	}
	urlID := data.ParsedURL.Path[1:] // Removes the first / from path
	longUrl, err := rs.cache.Get(r.Context(), urlID)
	if err != nil {
		rs.logger.Error("unable to retrieve long url from cache", "error", err)
		render.Render(w, r, ErrInternalServerError(err))
		return
	}
	if longUrl == "" {
		rs.logger.Warn("unable to find long url in cache", "error", err)
		render.Render(w, r, ErrNotFound(errors.New("long url not found")))
		return
	}
	rs.logger.Info("Decoded url",
		"shortUrl", data.URL,
		"urlId", urlID,
		"longUrl", longUrl,
	)
	render.Render(w, r, &URLPayload{URL: longUrl})
}
