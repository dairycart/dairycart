package postgres

import (
	"database/sql"

	"github.com/dairycart/dairycart/api/storage"
)

var _ storage.Storage = (*Postgres)(nil)

type Postgres struct{ *sql.DB }
