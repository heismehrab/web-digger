package services

import (
	"context"
	"web-digger/pkg/logger"
)

type AnalyzerService struct {
	logger *logger.StandardLogger
}

func NewAnalyzerService(logger *logger.StandardLogger) *AnalyzerService {
	return &AnalyzerService{
		logger: logger,
	}
}

func (a *AnalyzerService) Analyze(ctx context.Context, URL string) (any, error) {
	return 0, nil
}
