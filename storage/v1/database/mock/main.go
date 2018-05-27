package dairymock

import (
	"database/sql"

	"github.com/dairycart/dairycart/storage/v1/database"

	"github.com/spf13/viper"
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

func (m *MockDB) Migrate(db *sql.DB, cfg *viper.Viper) error {
	args := m.Called(db, cfg)

	return args.Error(0)
}
