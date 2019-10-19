package huh

import (
	"database/sql"

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
}

// Config establish a DB connection
func Config(dialect string, dbConfig DBConfig) {
	if dialect == "mysql" {
		db, err := sql.Open("mysql", dbConfig.Master)
		checkError(err)

		currentDB = &huhDB{*db}
		return
	}
	panic(ErrDialectNotSupported)
}
