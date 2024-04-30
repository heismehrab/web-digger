package ports

import "context"

type AnalyzerService interface {
	Analyze(ctx context.Context, URL string) (any, error)
}
