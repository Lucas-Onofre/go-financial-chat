package mocks

import "github.com/stretchr/testify/mock"

type MockBroker struct {
	mock.Mock
}

func (m *MockBroker) Publish(queue string, message string) error {
	args := m.Called(queue, message)
	return args.Error(0)
}

func (m *MockBroker) Subscribe(queue string, handler func(message string) error) error {
	args := m.Called(queue, handler)
	return args.Error(0)
}

func (m *MockBroker) Close() error {
	args := m.Called()
	return args.Error(0)
}
