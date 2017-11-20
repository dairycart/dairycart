package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
)

func (m *MockDB) GetLoginAttempt(id uint64) (models.LoginAttempt, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) CreateLoginAttempt(nu models.LoginAttempt) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0), args.Get(1), args.Error(2)
}

func (m *MockDB) UpdateLoginAttempt(updated models.LoginAttempt) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) DeleteLoginAttempt(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}
