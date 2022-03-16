package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/darioblanco/shortesturl/app/internal/cache"
	"github.com/darioblanco/shortesturl/app/internal/config"
	apphttp "github.com/darioblanco/shortesturl/app/internal/http"
	"github.com/darioblanco/shortesturl/app/internal/logging"
	_ "github.com/darioblanco/shortesturl/docs"
)

// An Application holds the configuration, context, logger models and router
// needed to run shortesturl
type Application interface {
	Serve()
}

type application struct {
	conf   *config.Values
	ctx    context.Context
	logger logging.Logger
	router http.Handler
}

// New creates an application instance
func New(ctx context.Context) Application {
	// Config
	configFilename := "config.dev"
	if _, err := os.Stat("cmd/config.dev.yaml"); errors.Is(err, os.ErrNotExist) {
		configFilename = "config.default"
	}
	conf, err := config.New(configFilename, "./cmd")
	if err != nil {
		log.Fatalf("Unable to load configuration: %v", err)
	}

	// Logger
	logger, err := logging.NewLogger(conf)
	if err != nil {
		log.Fatalf("Unable to load logger: %v", err)
	}
	logger.Info("Loaded config", "filename", configFilename)

	// Cache
	c, err := cache.New(ctx, conf, logger)
	if err != nil {
		log.Fatalf("Unable to connect to cache: %v", err)
	}

	// Router
	router, err := apphttp.NewRouter(
		ctx,
		conf,
		logger,
		c,
	)
	if err != nil {
		log.Fatalf("Unable to load router: %v", err)
	}

	return &application{
		conf:   conf,
		ctx:    ctx,
		logger: logger,
		router: router,
	}
}

// Serve sets the application ready to receive and process requests
func (a *application) Serve() {
	// The HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", a.conf.HttpHost, a.conf.HttpPort),
		Handler: a.router,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(a.ctx)

	// Listen for syscall signals for process to interrupt/quit
	// See https://github.com/go-chi/chi/blob/master/_examples/graceful/main.go
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatalf(
					"Graceful shutdown timed out... Forcing exit: %v",
					context.DeadlineExceeded,
				)
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatalf("Unable to shutdown server: %v", err)
		}
		serverStopCtx()
	}()
	// Run the server
	a.logger.Info("HTTP server started", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Unable to start HTTP server: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
}
