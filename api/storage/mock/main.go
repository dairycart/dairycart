package dairymock

import (
	"database/sql"

	"github.com/dairycart/dairycart/api/storage"

	"github.com/stretchr/testify/mock"
)

var _ storage.Storage = (*MockDB)(nil)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Begin() (*sql.Tx, error) {
	args := m.Called()

	return args.Get(0).(*sql.Tx), args.Error(1)
}
