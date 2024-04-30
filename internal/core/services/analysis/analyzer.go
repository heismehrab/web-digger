package analysis

import (
	"context"
	"golang.org/x/net/html"
	"io"
	"strings"
	"sync"
	"time"
	"web-digger/internal/core/domain/models"
	"web-digger/internal/core/ports"
	"web-digger/pkg/logger"
)

type AnalyzerService struct {
	logger     *logger.StandardLogger
	pageParser ports.ParserAdapter
	result     *models.AnalyzerResult

	sync.WaitGroup
	sync.Mutex
}

// NewAnalyzerService Return an instance of Analyzer service.
func NewAnalyzerService(logger *logger.StandardLogger, pageParser ports.ParserAdapter) *AnalyzerService {
	return &AnalyzerService{
		logger:     logger,
		pageParser: pageParser,
	}
}

// Analyze Tries to parse the page first and then inspect it to fetch required details.
func (a *AnalyzerService) Analyze(ctx context.Context, url string) (*models.AnalyzerResult, error) {
	allLinksAreProcessed := make(chan bool, 1)
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

	a.WaitGroup.Add(1)
	go a.getHTMLVersion(ctx, parsedHTMLPage)

	a.WaitGroup.Add(1)
	go a.getHTMLTitle(ctx, parsedHTMLPage)

	a.WaitGroup.Add(1)
	go a.getHTMLHeadings(ctx, parsedHTMLPage)

	a.WaitGroup.Add(1)
	go a.getHTMLALinks(ctx, parsedHTMLPage, url, allLinksAreProcessed)

	a.WaitGroup.Add(1)
	go a.getHTMLALinksStatus(ctx, allLinksAreProcessed)

	// Waiting for goroutines to finish.
	a.WaitGroup.Wait()

	a.logger.InfoF("all goroutines are done. | time: %d", time.Now().Unix())

	return a.result, nil
}
