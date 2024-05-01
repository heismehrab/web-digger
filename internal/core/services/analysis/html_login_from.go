package analysis

import (
	"context"
	"golang.org/x/net/html"
	"strings"
	"time"
)

// containsLoginForm Checks whether HTML content contains login form or not.
// For now, this functions detects forms that only have `action` key
// with value `login.
func (a *AnalyzerService) containsLoginForm(ctx context.Context, parsedHTMLPage string) {
	a.logger.InfoF("start getting HTML forms | time: %d", time.Now().Unix())

	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			// End of the HTML document
			break
		}

		if tokenType == html.StartTagToken {
			token := tokenizer.Token()

			if token.Data == "form" {
				// Check if the form has input fields for username and password.
				for _, attr := range token.Attr {
					if attr.Key == "action" && strings.Contains(attr.Val, "login") {
						a.result.HasLoginForm = true
						a.WaitGroup.Done()

						return
					}
				}
			}
		}
	}

	a.result.HasLoginForm = false
	a.WaitGroup.Done()
}
