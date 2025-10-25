package mocks

import "github.com/stretchr/testify/mock"

type MockMarketDataProvider struct {
	mock.Mock
}

func (m *MockMarketDataProvider) GetMarketData(stockCommand string) (string, error) {
	args := m.Called(stockCommand)
	return args.String(0), args.Error(1)
}
