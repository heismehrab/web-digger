package ports

import "context"

type ParserAdapter interface {
	ParseWebPage(ctx context.Context, url string) (string, error)
	GetWebPageStatus(ctx context.Context, url string) error
}
