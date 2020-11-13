package db

import (
	"testing"

	"github.com/golang-tire/pkg/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

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

	dbc, err := gorm.Open(postgres.Open(testDSN.String()), &gorm.Config{})
	if err != nil {
		log.Error("open database failed", log.Err(err))
		return nil
	}

	err = CreateSchema(dbc, models)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	return &DB{dbc}
}

// ResetTables truncates all data in the specified tables.
func ResetTables(t *testing.T, db *DB, tables ...string) error {
	return db.DB().Migrator().DropTable(tables)
}
