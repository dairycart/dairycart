package dairymock

import (
	"database/sql"

	"github.com/dairycart/dairycart/storage/database"

	"github.com/stretchr/testify/mock"
)

var _ database.Storer = (*MockDB)(nil)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() (*sql.Tx, error) {
	args := m.Called()

	return args.Get(0).(*sql.Tx), args.Error(1)
}
