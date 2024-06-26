package app

import (
	"context"
	"os"
	"sync"
	"web-digger/internal/app/handlers/http"
	"web-digger/internal/config"
	PageParser "web-digger/internal/core/infrastructure/http"
	"web-digger/internal/core/services/analysis"
	"web-digger/pkg/logger"
)

// Master Using sync.WaitGroup to ensure graceful
// shutdown waits for background jobs before to exit gracefully too.
type Master struct {
	sync.WaitGroup
	cfg         config.Config
	restHandler *http.Handler
	logger      *logger.StandardLogger
}

func NewInstance(cfg config.Config) *Master {
	return &Master{
		cfg: cfg,
	}
}

func (m *Master) Bootstrap(ctx context.Context, standardLogger *logger.StandardLogger) error {
	// Initiate services, dependencies, etc...

	// Create Logger instance.
	m.logger = standardLogger

	// Initializing Application services.
	analyzerService := analysis.NewAnalyzerService(
		m.logger,
		PageParser.NewPageParser(PageParser.WithHTTPTimeOut(PageParser.LongTimeout)),
	)

	// Create Application's HTTP handler.
	m.restHandler = http.CreateHandler(
		standardLogger,
		m.cfg,

		// Please define services here.
		analyzerService,
	)

	return nil
}

// Start is used to start the Rest handler.
func (m *Master) Start(ctx context.Context) {
	m.restHandler.StartBlocking(ctx, m.cfg.App.Port)
}

func (m *Master) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// Wait for OS signals
	<-quitSignal

	m.restHandler.Stop()

	close(done)
}
