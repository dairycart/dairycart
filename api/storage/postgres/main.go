package postgres

import (
	"github.com/dairycart/dairycart/api/storage"
)

var _ storage.Storer = (*postgres)(nil)

type postgres struct{}

func NewPostgres() *postgres {
	return &postgres{}
}
