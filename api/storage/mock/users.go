package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/api/storage"
	"github.com/dairycart/dairycart/api/storage/models"
)

func (m *MockDB) UserExists(db storage.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetUser(db storage.Querier, id uint64) (*models.User, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetUserList(db storage.Querier, qf *models.QueryFilter) ([]models.User, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockDB) GetUserCount(db storage.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateUser(db storage.Querier, nu *models.User) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateUser(db storage.Querier, updated *models.User) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteUser(db storage.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
