package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) GetLoginAttempt(id uint64) (*models.LoginAttempt, error) {
	args := m.Called(id)
	return args.Get(0).(*models.LoginAttempt), args.Error(1)
}

func (m *MockDB) CreateLoginAttempt(nu *models.LoginAttempt) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateLoginAttempt(updated *models.LoginAttempt) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteLoginAttempt(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0).(time.Time), args.Error(1)
}
