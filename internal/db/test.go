package db

import (
	"fmt"
	"testing"

	"github.com/go-pg/pg/v10"
	"github.com/golang-tire/pkg/config"
)

var testDB *DB

var (
	testDSN config.String
)

// NewForTest returns the database connection for testing purpose.
func NewForTest(t *testing.T, models []interface{}) *DB {
	if testDB != nil {
		return testDB
	}

	testDSN = config.RegisterString("db.testDSN", "postgres://postgres:postgres@localhost:5432/auth_test?sslmode=disable")
	opt, err := pg.ParseURL(testDSN.String())
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	dbc := pg.Connect(opt)

	err = CreateSchema(dbc, models)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return &DB{dbc}
}

// ResetTables truncates all data in the specified tables.
func ResetTables(t *testing.T, db *DB, tables ...string) {
	for _, table := range tables {
		_, err := db.DB().Exec(fmt.Sprintf("TRUNCATE %s", table))
		if err != nil {
			t.Error(err)
			t.FailNow()
		}
	}
}
