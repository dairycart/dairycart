package postgres

import (
	"database/sql"
)

type Postgres struct {
	DB *sql.DB
}
