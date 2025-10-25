package mocks

import (
	"gorm.io/gorm"

	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Create(entity any) *gorm.DB {
	args := m.Called(entity)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Where(query any, args ...any) *gorm.DB {
	calledArgs := m.Called(append([]any{query}, args...)...)
	return calledArgs.Get(0).(*gorm.DB)
}

func (m *MockDB) First(dest any, conds ...any) *gorm.DB {
	calledArgs := m.Called(append([]any{dest}, conds...)...)
	return calledArgs.Get(0).(*gorm.DB)
}
