package dairymock

import (
	"time"

	"github.com/dairycart/dairycart/models/v1"
	"github.com/dairycart/dairycart/storage/v1/database"
)

func (m *MockDB) GetUserByUsername(db database.Querier, username string) (*models.User, error) {
	args := m.Called(db, username)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) UserWithUsernameExists(db database.Querier, username string) (bool, error) {
	args := m.Called(db, username)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) UserExists(db database.Querier, id uint64) (bool, error) {
	args := m.Called(db, id)
	return args.Bool(0), args.Error(1)
}

func (m *MockDB) GetUser(db database.Querier, id uint64) (*models.User, error) {
	args := m.Called(db, id)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockDB) GetUserList(db database.Querier, qf *models.QueryFilter) ([]models.User, error) {
	args := m.Called(db, qf)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockDB) GetUserCount(db database.Querier, qf *models.QueryFilter) (uint64, error) {
	args := m.Called(db, qf)
	return args.Get(0).(uint64), args.Error(1)
}

func (m *MockDB) CreateUser(db database.Querier, nu *models.User) (uint64, time.Time, error) {
	args := m.Called(db, nu)
	return args.Get(0).(uint64), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockDB) UpdateUser(db database.Querier, updated *models.User) (time.Time, error) {
	args := m.Called(db, updated)
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockDB) DeleteUser(db database.Querier, id uint64) (time.Time, error) {
	args := m.Called(db, id)
	return args.Get(0).(time.Time), args.Error(1)
}
