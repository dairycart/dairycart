package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairymodels/v1"
)

func (m *MockDB) LoginAttemptsHaveBeenExhausted(db storage.Querier, username string) (bool, error) {
	args := m.Called(db, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) LoginAttemptExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetLoginAttempt(db storage.Querier, id uint64) (*models.LoginAttempt, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.LoginAttempt), args.Error(1)
}

func (m *MockDB) GetLoginAttemptList(db storage.Querier, qf *models.QueryFilter) ([]models.LoginAttempt, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.LoginAttempt), args.Error(1)
}

func (m *MockDB) GetLoginAttemptCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateLoginAttempt(db storage.Querier, nu *models.LoginAttempt) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateLoginAttempt(db storage.Querier, updated *models.LoginAttempt) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteLoginAttempt(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
