package analysis

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
	"web-digger/internal/core/domain/models"
	"web-digger/pkg/logger"
	"web-digger/pkg/mocks"
)

const StubSampleValidUrl = "https://www.google.com"

var config = logger.CreateLogger(logger.Config{
	LogLevel:      slog.LevelDebug,
	GrayLogActive: false,
})

func TestGetHTMLVersion(t *testing.T) {
	doc := `<!invalid doctype>`

	mockedParser := new(mocks.MockedParser)
	mockedParser.On("ParseWebPage", mock.Anything, mock.Anything).Return(
		doc,
		nil,
	)

	expected := ``
	analyzerService := NewAnalyzerService(config, mockedParser)
	analyzerService.result = &models.AnalyzerResult{}

	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE HTML>`
	expected = `HTML 5`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN">`
	expected = `HTML 4.01 Strict`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">`
	expected = `HTML 4.01 Transitional`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01 Frameset//EN">`
	expected = `HTML 4.01 Frameset`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN">`
	expected = `XHTML 1.0 Strict`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN">`
	expected = `XHTML 1.0 Transitional`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Frameset//EN">`
	expected = `XHTML 1.0 Frameset`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)

	doc = `<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.1//EN">`
	expected = `XHTML 1.1`
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLVersion(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Version)
}

func TestGetHTMLTitle(t *testing.T) {
	doc := `<title>title-1</title>`

	mockedParser := new(mocks.MockedParser)
	mockedParser.On("ParseWebPage", mock.Anything, mock.Anything).Return(
		doc,
		nil,
	)

	expected := `title-1`
	analyzerService := NewAnalyzerService(config, mockedParser)
	analyzerService.result = &models.AnalyzerResult{}
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLTitle(context.Background(), doc)

	assert.Equal(t, expected, analyzerService.result.Title)
}

func TestGetHTMLHTags(t *testing.T) {
	doc := `<h1>1</h1><h1>2</h1><h2>1</h2><h3>2</h3><h6>1</h6><h6>2</h6><h5>1</h5><h6>3</h6>`

	mockedParser := new(mocks.MockedParser)
	mockedParser.On("ParseWebPage", mock.Anything, mock.Anything).Return(
		doc,
		nil,
	)

	analyzerService := NewAnalyzerService(config, mockedParser)
	analyzerService.result = &models.AnalyzerResult{}
	analyzerService.WaitGroup.Add(1)
	analyzerService.getHTMLHeadings(context.Background(), doc)

	assert.Equal(t, 2, analyzerService.result.Hs[1])
	assert.Equal(t, 1, analyzerService.result.Hs[2])
	assert.Equal(t, 1, analyzerService.result.Hs[3])
	assert.Equal(t, 3, analyzerService.result.Hs[6])
	assert.Equal(t, 1, analyzerService.result.Hs[5])
}

func TestGetHTMLALinks(t *testing.T) {
	internalURL := "https://internall.com"

	doc := `
<a href="%s/url-1.com">1</a>
<a href="%s/url-1.com">2</a>
<a href="#">0</a>
<a href="https://external.com">3</a>
<a href="https://external.com">4</a>
`

	mockedParser := new(mocks.MockedParser)
	mockedParser.On("ParseWebPage", mock.Anything, mock.Anything).Return(
		doc,
		nil,
	)

	analyzerService := NewAnalyzerService(config, mockedParser)
	analyzerService.result = &models.AnalyzerResult{}
	analyzerService.WaitGroup.Add(1)

	analyzerService.getHTMLALinks(
		context.Background(),
		fmt.Sprintf(doc, internalURL, internalURL),
		internalURL,
		make(chan bool),
	)

	assert.Equal(t, 2, len(analyzerService.result.InternalLinks))
	assert.Equal(t, 2, len(analyzerService.result.ExternalLinks))
}

func TestHTMHasLoginForm(t *testing.T) {
	doc := `<form action="login">`

	mockedParser := new(mocks.MockedParser)
	mockedParser.On("ParseWebPage", mock.Anything, mock.Anything).Return(
		doc,
		nil,
	)

	analyzerService := NewAnalyzerService(config, mockedParser)
	analyzerService.result = &models.AnalyzerResult{}
	analyzerService.WaitGroup.Add(1)

	analyzerService.containsLoginForm(
		context.Background(),
		doc,
	)

	assert.True(t, analyzerService.result.HasLoginForm)

	// Test without targeted action attribute.
	doc2 := `<form action="">`
	analyzerService.result = &models.AnalyzerResult{}
	analyzerService.WaitGroup.Add(1)

	analyzerService.containsLoginForm(
		context.Background(),
		doc2,
	)

	assert.False(t, analyzerService.result.HasLoginForm)
}
