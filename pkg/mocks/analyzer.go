package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
	"web-digger/internal/core/domain/models"
)

type MockedAnalyzerService struct {
	mock.Mock
}

func (m *MockedAnalyzerService) Analyze(ctx context.Context, URL string) (*models.AnalyzerResult, error) {
	args := m.Called(ctx, URL)

	return args.Get(0).(*models.AnalyzerResult), args.Error(1)
}
