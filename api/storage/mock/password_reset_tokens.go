package dairymock

import (
	"database/sql"
	"time"

	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) PasswordResetTokenExists(id uint64) (bool, error) {
	args := m.Called(id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetPasswordResetToken(id uint64) (*models.PasswordResetToken, error) {
	args := m.Called(id)
	return args.Get(0).(*models.PasswordResetToken), args.Error(1)
}

func (m *MockDB) CreatePasswordResetToken(nu *models.PasswordResetToken) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdatePasswordResetToken(updated *models.PasswordResetToken) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeletePasswordResetToken(id uint64, tx *sql.Tx) (time.Time, error) {
	args := m.Called(id, tx)
	return args.Get(0).(time.Time), args.Error(1)
}
