package mocks

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockedParser struct {
	mock.Mock
}

func (m *MockedParser) ParseWebPage(ctx context.Context, url string) (string, error) {
	args := m.Called(ctx, url)

	return args.String(0), args.Error(1)
}

func (m *MockedParser) GetWebPageStatus(ctx context.Context, url string) error {
	args := m.Called(ctx, url)

	return args.Error(0)
}
