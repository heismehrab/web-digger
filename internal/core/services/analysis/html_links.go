package analysis

import (
	"context"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"time"
)

// getHTMLVersion finds the title of HTML page
// via finding html.StartTagToken in given tokens.
func (a *AnalyzerService) getHTMLALinks(
	ctx context.Context,
	parsedHTMLPage string,
	webPageUrl string,
	allLinksAreProcessed chan<- bool,
) {
	a.logger.InfoF("start getting HTML A tags | time: %d", time.Now().Unix())

	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	for {
		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			allLinksAreProcessed <- true

			a.WaitGroup.Done()

			return
		}

		token := tokenizer.Token()

		if tt == html.StartTagToken && token.Data == "a" {

			for _, attr := range token.Attr {

				if attr.Key != "href" {
					continue
				}

				// empty or # urls must be ignored.
				if attr.Val == "#" {
					continue
				}

				// Parsing URL.
				parsedURL, err := url.Parse(strings.TrimSpace(attr.Val))

				if err != nil {
					a.logger.With("url", attr.Val).ErrorContext(ctx, err.Error())

					continue
				}

				if !strings.HasPrefix(parsedURL.Scheme, "http") {
					continue
				}

				if strings.Contains(webPageUrl, strings.ToLower(parsedURL.Host)) {
					a.result.InternalLinks = append(a.result.InternalLinks, parsedURL.String())

					continue
				}

				a.result.ExternalLinks = append(a.result.ExternalLinks, parsedURL.String())
			}
		}
	}

	allLinksAreProcessed <- true
	a.WaitGroup.Done()
}
