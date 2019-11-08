package huh

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type huhDB struct {
	sql.DB
}

var currentDB *huhDB

// DBConfig stores db connecting addresses
// Slaves is not supported for now
type DBConfig struct {
	Master string
	Slaves []string
	Option DBOption
}

// DBOption db options
type DBOption struct {
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// Config establish a DB connection
func Config(dialect string, dbConfig DBConfig) {
	if dialect == "mysql" {
		db, err := sql.Open("mysql", dbConfig.Master)
		checkError(err)

		// set sql.DB options
		if &dbConfig.Option.ConnMaxLifetime != nil {
			db.SetConnMaxLifetime(dbConfig.Option.ConnMaxLifetime)
		}
		if &dbConfig.Option.MaxOpenConns != nil {
			db.SetMaxOpenConns(dbConfig.Option.MaxOpenConns)
		}
		if &dbConfig.Option.MaxIdleConns != nil {
			db.SetMaxIdleConns(dbConfig.Option.MaxIdleConns)
		}

		currentDB = &huhDB{*db}
		return
	}
	panic(ErrDialectNotSupported)
}

// Close DB connection
func Close() error {
	return currentDB.Close()
}
