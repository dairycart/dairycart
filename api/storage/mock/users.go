package postgres

import (
	"time"

	"github.com/verygoodsoftwarenotvirus/gnorm-dairymodels/models"
)

func (m *MockDB) GetUser(id uint64) (models.User, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) CreateUser(nu models.User) (uint64, time.Time, error) {
	args := m.Called(nu)
	return args.Get(0), args.Get(1), args.Error(2)
}

func (m *MockDB) UpdateUser(updated models.User) (time.Time, error) {
	args := m.Called(updated)
	return args.Get(0), args.Error(1)
}

func (m *MockDB) DeleteUser(id uint64) (time.Time, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}
