package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) LoginAttemptsHaveBeenExhausted(db database.Querier, username string) (bool, error) {
	args := m.Called(db, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) LoginAttemptExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetLoginAttempt(db database.Querier, id uint64) (*models.LoginAttempt, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.LoginAttempt), args.Error(1)
}

func (m *MockDB) GetLoginAttemptList(db database.Querier, qf *models.QueryFilter) ([]models.LoginAttempt, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.LoginAttempt), args.Error(1)
}

func (m *MockDB) GetLoginAttemptCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateLoginAttempt(db database.Querier, nu *models.LoginAttempt) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateLoginAttempt(db database.Querier, updated *models.LoginAttempt) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteLoginAttempt(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
