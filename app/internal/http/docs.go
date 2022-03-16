package http

import (
	"fmt"
	"net/url"
	"path"

	"github.com/darioblanco/shortesturl/app/internal/config"
	"github.com/go-chi/chi/v5"
	swagger "github.com/swaggo/http-swagger"
)

type docs struct {
	config *config.Values
}

// @title ShortestURL API
// @version 1.0
// @description ShortestURL API.
// @termsOfService http://swagger.io/terms/

// @contact.name Darío Blanco Iturriaga
// @contact.url https://darioblanco.com
// @contact.email dblancoit@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost
// @BasePath /
// @query.collection.format multi
func (rs docs) Router() chi.Router {
	r := chi.NewRouter()
	swaggerJSONUrl := url.URL{
		Scheme: rs.config.HttpScheme,
		Host:   fmt.Sprintf("%s:%d", rs.config.HttpHost, rs.config.HttpPort),
		Path:   path.Join("docs", "doc.json"),
	}
	r.Get("/*", swagger.Handler(
		swagger.URL(swaggerJSONUrl.String()),
	))
	return r
}
