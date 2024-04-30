package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"testing"
	"web-digger/internal/config"
	"web-digger/pkg/logger"
)

func setupHandler(t *testing.T) *Handler {
	consoleLogger := logger.CreateLogger(logger.Config{
		LogLevel:      slog.LevelDebug,
		GrayLogActive: false,
	})

	return CreateHandler(consoleLogger, config.Config{})
}

func createRequest(t *testing.T, method string, url string, body interface{}, headers map[string]string) *http.Request {
	var reader io.Reader

	if body != nil {
		marshaledBody, err := json.Marshal(body)

		if err != nil {
			t.Fatalf("failed to marshal request body with error: %s", err.Error())
		}

		reader = bytes.NewBuffer(marshaledBody)
	}

	req, err := http.NewRequestWithContext(
		context.Background(),
		method,
		url,
		reader,
	)

	if err != nil {
		t.Fatalf("failed to create request with error: %s", err.Error())
	}

	return req
}
