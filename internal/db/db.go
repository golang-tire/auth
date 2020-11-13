package db

import (
	"context"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-tire/pkg/config"
	"github.com/golang-tire/pkg/log"
)

// DB represents a DB connection that can be used to run SQL queries.
type DB struct {
	db *gorm.DB
}

type contextKey int

const (
	txKey contextKey = iota
)

var (
	host     config.String
	port     config.Int
	db       config.String
	password config.String
	user     config.String
)

// DB returns the dbx.DB wrapped by this object.
func (db *DB) DB() *gorm.DB {
	return db.db
}

// With returns a Builder that can be used to build and execute SQL queries.
// With will return the transaction if it is found in the given context.
// Otherwise it will return a DB connection associated with the context.
func (db *DB) With(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey).(*gorm.DB); ok {
		return tx
	}
	return db.db.WithContext(ctx)
}

func Init(ctx context.Context) (*DB, error) {

	host = config.RegisterString("db.host", "localhost")
	port = config.RegisterInt("db.port", 5432)
	db = config.RegisterString("db.db", "postgres")
	password = config.RegisterString("db.password", "postgres")
	user = config.RegisterString("db.user", "postgres")

	if err := config.Load(); err != nil {
		log.Panic("load database settings failed")
		return nil, err
	}

	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable TimeZone=Europe/Stockholm",
		user.String(), password.String(), db.String(), host.String(), port.Int(),
	)
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error("open database failed", log.Err(err))
		return nil, err
	}

	go func() {
		<-ctx.Done()
		pgDB, err := database.DB()
		if err != nil {
			log.Error("error in close database", log.Err(err))
		}
		if err := pgDB.Close(); err != nil {
			log.Error("error in close database", log.Err(err))
		}
	}()

	return &DB{database}, nil
}

// CreateSchema creates database schema for test
func CreateSchema(db *gorm.DB, models []interface{}) error {
	return db.AutoMigrate(models...)
}
