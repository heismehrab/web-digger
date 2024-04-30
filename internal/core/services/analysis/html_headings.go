package analysis

import (
	"context"
	"golang.org/x/net/html"
	"strings"
	"time"
	"web-digger/internal/core/domain/models"
)

// getHTMLHeadings finds the h tags of HTML page
// via finding html.StartTagToken
func (a *AnalyzerService) getHTMLHeadings(ctx context.Context, parsedHTMLPage string) {
	a.logger.InfoF("start getting HTML H tags | time: %d", time.Now().Unix())

	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	tags := models.HTags{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
	}

	for {
		tt := tokenizer.Next()

		if tt == html.ErrorToken {
			break
		}

		token := tokenizer.Token().Data

		if tt == html.StartTagToken && string(token[0]) == "h" {
			switch token {
			case "h1":
				tags[1]++
			case "h2":
				tags[2]++
			case "h3":
				tags[3]++
			case "h4":
				tags[4]++
			case "h5":
				tags[5]++
			case "h6":
				tags[6]++
			}
		}
	}

	a.result.Hs = tags
	a.WaitGroup.Done()
}
