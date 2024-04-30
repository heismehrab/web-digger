package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
	"web-digger/internal/config"
	"web-digger/pkg/logger"
)

const (
	DefaultTimeOutForGracefulShutDown = 5 * time.Second
	IdleTimeout                       = time.Second * 60
	ReadTimeout                       = time.Second * 15
	WriteTimeout                      = time.Second * 15
	MaxMultipartMemory                = 8 << 20 // 8 MiB
)

type Handler struct {
	// HTTPServer is the main server object
	HTTPServer *http.Server

	config config.Config
	logger *logger.StandardLogger
}

// CreateHandler Creates a new instance of REST handler.
func CreateHandler(
	logger *logger.StandardLogger,
	config config.Config,
) *Handler {
	return &Handler{
		logger: logger,
		config: config,
	}
}

// StartBlocking starts the http server.
func (h *Handler) StartBlocking(ctx context.Context, defaultPort int) {
	addr := fmt.Sprintf(":%v", defaultPort)

	h.logger.InfoF("[OK] Starting HTTP REST Server on %s ", addr)

	h.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: h.setupRouter(),

		// Good practice to set timeouts to avoid Slow-loris attacks.
		WriteTimeout: WriteTimeout,
		ReadTimeout:  ReadTimeout,
		IdleTimeout:  IdleTimeout,
	}

	err := h.HTTPServer.ListenAndServe()

	if !errors.Is(err, http.ErrServerClosed) {
		h.logger.Error(err.Error())
	}

	// Code Reach Here after HTTP Server Shutdown!
	h.logger.Info("[OK] HTTP REST Server is shutting down!")
}

// Stop handles the http server in graceful shutdown.
func (h *Handler) Stop() {
	// Create a 5s timeout context or waiting for app to shut down after 5 seconds.
	ctxTimeout, cancelTimeout := context.WithTimeout(
		context.Background(),
		DefaultTimeOutForGracefulShutDown,
	)

	defer cancelTimeout()

	h.HTTPServer.SetKeepAlivesEnabled(false)

	if err := h.HTTPServer.Shutdown(ctxTimeout); err != nil {
		h.logger.Error(err.Error())
	}

	h.logger.Info("[OK] HTTP REST Server graceful shutdown completed")
}
