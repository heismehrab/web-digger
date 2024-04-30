package analysis

import (
	"context"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
	"web-digger/internal/core/domain/models"
	"web-digger/internal/core/infrastructure/http"
	"web-digger/pkg/logger"
)

var ErrHTMLVersionNotFound = errors.New("failed to get HTML version")

type AnalyzerService struct {
	logger     *logger.StandardLogger
	pageParser *http.PageParser
	result     *models.AnalyzerResult
}

// NewAnalyzerService Return an instance of Analyzer service.
func NewAnalyzerService(logger *logger.StandardLogger, pageParser *http.PageParser) *AnalyzerService {
	return &AnalyzerService{
		logger:     logger,
		pageParser: pageParser,
	}
}

// Analyze Tries to parse the page first and then inspect it to fetch required details.
func (a *AnalyzerService) Analyze(ctx context.Context, url string) (*models.AnalyzerResult, error) {
	a.result = &models.AnalyzerResult{}

	// Parse URL and fetch its contents.
	parsedHTMLPage, err := a.pageParser.ParseWebPage(ctx, url)

	if err != nil {
		a.logger.With("url", url).ErrorContext(ctx, err.Error())

		return nil, err
	}

	// Validate html content.
	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	for {
		t := tokenizer.Next()

		if t == html.ErrorToken {
			err := tokenizer.Err()

			if err == io.EOF {
				break
			}

			a.logger.With("url", url).ErrorContext(
				ctx,
				tokenizer.Err().Error(),
			)

			return nil, err
		}
	}

	a.getHTMLVersion(ctx, parsedHTMLPage)
	a.getHTMLTitle(ctx, parsedHTMLPage)
	a.getHTMLHeadings(ctx, parsedHTMLPage)

	return a.result, nil
}
