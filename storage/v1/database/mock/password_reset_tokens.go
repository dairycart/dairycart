package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) PasswordResetTokenForUserIDExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) PasswordResetTokenWithTokenExists(db database.Querier, token string) (bool, error) {
	args := m.Called(db, token)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) PasswordResetTokenExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetPasswordResetToken(db database.Querier, id uint64) (*models.PasswordResetToken, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.PasswordResetToken), args.Error(1)
}

func (m *MockDB) GetPasswordResetTokenList(db database.Querier, qf *models.QueryFilter) ([]models.PasswordResetToken, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.PasswordResetToken), args.Error(1)
}

func (m *MockDB) GetPasswordResetTokenCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreatePasswordResetToken(db database.Querier, nu *models.PasswordResetToken) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdatePasswordResetToken(db database.Querier, updated *models.PasswordResetToken) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeletePasswordResetToken(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
