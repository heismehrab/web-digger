package analysis

import (
	"context"
	"golang.org/x/net/html"
	"strings"
	"time"
)

// getHTMLTitle finds the title of HTML page
// via finding html.StartTagToken and html.TextToken in given tokens.
func (a *AnalyzerService) getHTMLTitle(ctx context.Context, parsedHTMLPage string) {
	a.logger.InfoF("start getting HTML title | time: %d", time.Now().Unix())

	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	for {
		tt := tokenizer.Next()
		token := tokenizer.Token().Data

		if tt == html.ErrorToken {
			a.result.Title = "failed to fetch page title"
			a.WaitGroup.Done()

			return
		}

		if tt == html.StartTagToken && token == "title" {
			tt = tokenizer.Next()

			if tt == html.TextToken {
				data := tokenizer.Token().Data

				if data != "" {
					a.result.Title = data
					a.WaitGroup.Done()

					return
				}

			}
		}
	}
}
