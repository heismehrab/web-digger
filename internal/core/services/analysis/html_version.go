package analysis

import (
	"context"
	"golang.org/x/net/html"
	"strings"
)

type htmlVersions struct {
	version, matcher string
}

// getHTMLVersion Indicates the version of HTML page
// via finding html.DoctypeToken in given tokens.
func (a *AnalyzerService) getHTMLVersion(ctx context.Context, parsedHTMLPage string) {
	tokenizer := html.NewTokenizer(strings.NewReader(parsedHTMLPage))

	for {
		tt := tokenizer.Next()

		if tt != html.DoctypeToken {
			continue
		}

		docTypes := [8]htmlVersions{
			{version: "HTML 4.01 Strict", matcher: `"-//W3C//DTD HTML 4.01//EN"`},
			{version: "HTML 4.01 Transitional", matcher: `"-//W3C//DTD HTML 4.01 TRANSITIONAL//EN"`},
			{version: "HTML 4.01 Frameset", matcher: `"-//W3C//DTD HTML 4.01 FRAMESET//EN"`},
			{version: "XHTML 1.0 Strict", matcher: `"-//W3C//DTD XHTML 1.0 STRICT//EN`},
			{version: "XHTML 1.0 Transitional", matcher: `"-//W3C//DTD XHTML 1.0 TRANSITIONAL//EN"`},
			{version: "XHTML 1.0 Frameset", matcher: `"-//W3C//DTD XHTML 1.0 FRAMESET//EN"`},
			{version: "XHTML 1.1", matcher: `"-//W3C//DTD XHTML 1.1//EN"`},
			{version: "HTML 5", matcher: `HTML`},
		}

		v := tokenizer.Token().Data

		for _, d := range docTypes {
			ok := strings.Contains(strings.ToUpper(v), d.matcher)

			if ok {
				a.result.Version = d.version

				return
			}
		}
	}
}
