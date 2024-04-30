package analysis

import (
	"context"
	"golang.org/x/net/html"
	"strings"
)

// getHTMLVersion finds the title of HTML page
// via finding html.StartTagToken and html.TextToken in given tokens.
func (a *AnalyzerService) getHTMLTitle(ctx context.Context, parsedHTMLPage string) {
	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	for {
		tt := tokenizer.Next()
		td := tokenizer.Token().Data

		if tt == html.StartTagToken && td == "title" {
			tt = tokenizer.Next()

			if tt == html.TextToken {
				data := tokenizer.Token().Data

				if data != "" {
					a.result.Title = data
					break
				}

			}

			break
		}
	}
}
