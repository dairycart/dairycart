package postgres

import (
	"github.com/dairycart/dairycart/api/storage"
)

var _ storage.Storer = (*Postgres)(nil)

type Postgres struct{}

func NewPostgres() *Postgres {
	return &Postgres{}
}
