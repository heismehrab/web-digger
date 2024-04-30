package ports

import (
	"context"
	"web-digger/internal/core/domain/models"
)

type AnalyzerService interface {
	Analyze(ctx context.Context, URL string) (*models.AnalyzerResult, error)
}
