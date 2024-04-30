package http

import (
	"context"
	"errors"
	"io"
	"net/http"
	"strings"
)

const contentType = "text/html"

var errURLFetchFailed = errors.New("failed to fetch URL")
var errInvalidContentType = errors.New("content type must be 'text/html'")

type PageParser struct {
	httpClient *http.Client
}

// NewPageParser returns an instance of HTTP parser Client.
func NewPageParser(options ...Option) *PageParser {
	defaultHTTPClient := getDefaultHTTPClient()

	service := &PageParser{
		httpClient: defaultHTTPClient,
	}

	// Apply all the functional options to configure the client.
	for _, opt := range options {
		opt(service)
	}

	return service
}

// ParseWebPage Tries to fetch given url content and validate it.
// The result's content will be returned if url has suitable conditions.
func (p *PageParser) ParseWebPage(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {
		return "", err
	}

	res, err := p.httpClient.Do(req)

	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	// Check result's status code.
	if res.StatusCode != http.StatusOK {
		return "", errURLFetchFailed
	}

	// Validate URL content type.
	if !strings.HasPrefix(res.Header.Get("Content-Type"), contentType) {
		return "", errInvalidContentType
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	return string(body), nil
}
