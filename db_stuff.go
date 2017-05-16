package main

import (
	"context"
	"io"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// Database is a generic interface that should allow us to mock Database queries in tests
// Note that these functions are almost entirely copied from go-pg's DB struct, because
// that should be the primary implementer of this interface
type Database interface {
	// https://godoc.org/github.com/go-pg/pg#DB
	Begin() (*pg.Tx, error)
	Close() error
	Context() context.Context
	CopyFrom(reader io.Reader, query interface{}, params ...interface{}) (orm.Result, error)
	CopyTo(writer io.Writer, query interface{}, params ...interface{}) (orm.Result, error)
	CreateTable(model interface{}, opt *orm.CreateTableOptions) error
	Delete(model interface{}) error
	Exec(query interface{}, params ...interface{}) (res orm.Result, err error)
	ExecOne(query interface{}, params ...interface{}) (orm.Result, error)
	FormatQuery(dst []byte, query string, params ...interface{}) []byte
	Insert(model ...interface{}) error
	Listen(channels ...string) *pg.Listener
	Model(model ...interface{}) *orm.Query
	OnQueryProcessed(fn func(*pg.QueryProcessedEvent))
	Options() *pg.Options
	Prepare(q string) (*pg.Stmt, error)
	Query(model, query interface{}, params ...interface{}) (res orm.Result, err error)
	QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error)
	RunInTransaction(fn func(*pg.Tx) error) error
	Select(model interface{}) error
	String() string
	Update(model interface{}) error
	WithContext(ctx context.Context) *pg.DB
	WithParam(param string, value interface{}) *pg.DB
	WithTimeout(d time.Duration) *pg.DB
}
