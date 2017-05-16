package main

import (
	"context"
	"io"
	"time"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
)

// type MockDBCall struct {
// 	FunctionName string
// }

// Purely a mock-style object for testing
type MockDB struct {
	CallList []string
}

func NewMockDB() *MockDB {
	return &MockDB{}
}

func (db *MockDB) Insert(model ...interface{}) error {
	db.CallList = append(db.CallList, "Insert")
	// structType := reflect.ValueOf(model).Type().Name()
	return nil
}
func (db *MockDB) Model(model ...interface{}) *orm.Query {
	db.CallList = append(db.CallList, "Model")
	return orm.NewQuery(db, model...)
}
func (db *MockDB) Select(model interface{}) error {
	db.CallList = append(db.CallList, "Select")
	return nil
}
func (db *MockDB) Update(model interface{}) error {
	db.CallList = append(db.CallList, "Update")
	return nil
}

func (db *MockDB) RunInTransaction(fn func(*pg.Tx) error) error {
	db.CallList = append(db.CallList, "RunInTransaction")
	return nil
}

// The following methods should never need to be called. Maybe I can delete them and simply
// edit the interface?
func (db *MockDB) Begin() (*pg.Tx, error) {
	db.CallList = append(db.CallList, "Begin")
	return nil, nil
}

func (db *MockDB) Close() error {
	db.CallList = append(db.CallList, "Close")
	return nil
}

func (db *MockDB) Context() context.Context {
	db.CallList = append(db.CallList, "Context")
	return nil
}

func (db *MockDB) CopyFrom(reader io.Reader, query interface{}, params ...interface{}) (orm.Result, error) {
	db.CallList = append(db.CallList, "CopyFrom")
	return nil, nil
}

func (db *MockDB) CopyTo(writer io.Writer, query interface{}, params ...interface{}) (orm.Result, error) {
	db.CallList = append(db.CallList, "CopyTo")
	return nil, nil
}

func (db *MockDB) CreateTable(model interface{}, opt *orm.CreateTableOptions) error {
	db.CallList = append(db.CallList, "CreateTable")
	return nil
}

func (db *MockDB) Delete(model interface{}) error {
	db.CallList = append(db.CallList, "Delete")
	return nil
}

func (db *MockDB) Exec(query interface{}, params ...interface{}) (res orm.Result, err error) {
	db.CallList = append(db.CallList, "Exec")
	return nil, nil
}

func (db *MockDB) ExecOne(query interface{}, params ...interface{}) (orm.Result, error) {
	db.CallList = append(db.CallList, "ExecOne")
	return nil, nil
}

func (db *MockDB) FormatQuery(dst []byte, query string, params ...interface{}) []byte {
	db.CallList = append(db.CallList, "FormatQuery")
	return nil
}

func (db *MockDB) Listen(channels ...string) *pg.Listener {
	db.CallList = append(db.CallList, "Listen")
	return nil
}

func (db *MockDB) OnQueryProcessed(fn func(*pg.QueryProcessedEvent)) {
	db.CallList = append(db.CallList, "OnQueryProcessed")
}

func (db *MockDB) Options() *pg.Options {
	db.CallList = append(db.CallList, "Options")
	return nil
}

func (db *MockDB) Prepare(q string) (*pg.Stmt, error) {
	db.CallList = append(db.CallList, "Prepare")
	return nil, nil
}

func (db *MockDB) Query(model, query interface{}, params ...interface{}) (res orm.Result, err error) {
	db.CallList = append(db.CallList, "Query")
	return nil, nil
}

func (db *MockDB) QueryOne(model, query interface{}, params ...interface{}) (orm.Result, error) {
	db.CallList = append(db.CallList, "QueryOne")
	return nil, nil
}
func (db *MockDB) String() string {
	db.CallList = append(db.CallList, "String")
	return ""
}

func (db *MockDB) WithContext(ctx context.Context) *pg.DB {
	db.CallList = append(db.CallList, "WithContext")
	return nil
}

func (db *MockDB) WithParam(param string, value interface{}) *pg.DB {
	db.CallList = append(db.CallList, "WithParam")
	return nil
}

func (db *MockDB) WithTimeout(d time.Duration) *pg.DB {
	db.CallList = append(db.CallList, "WithTimeout")
	return nil
}
