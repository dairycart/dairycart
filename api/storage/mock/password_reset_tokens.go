package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
)

func (m *MockDB) GetPasswordResetToken(id uint64) (models.PasswordResetToken, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) CreatePasswordResetToken(nu models.PasswordResetToken) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0), args.Get(1), args.Error(2)
}

func (m *MockDB) UpdatePasswordResetToken(updated models.PasswordResetToken) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) DeletePasswordResetToken(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}
