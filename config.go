package huh

import (
	"database/sql"
	"reflect"
	"time"

	// mysql package
	_ "github.com/go-sql-driver/mysql"
)

type huhDB struct {
	sql.DB
}

type Pool struct {
	masterDB *huhDB
	slaveDBs []*huhDB
	slaveIdx int
}

var pool Pool

var masterDB *huhDB
var slaveDBs []*huhDB

// DBConfig stores db connecting addresses
// Slaves is not supported for now
type DBConfig struct {
	Master string
	Slaves []string
}

func SetConnMaxLifetime(d time.Duration) {
	callForAllDBs("SetConnMaxLifetime", d)
}

func SetMaxOpenConns(i int) {
	callForAllDBs("SetMaxOpenConns", i)
}

func SetMaxIdleConns(i int) {
	callForAllDBs("SetMaxIdleConns", i)
}

// call functions for all master and slave DBS
func callForAllDBs(methodName string, args ...interface{}) {
	in := make([]reflect.Value, len(args))
	for _, arg := range args {
		in = append(in, reflect.ValueOf(arg))
	}

	// masterDB call
	masterDBValue := reflect.ValueOf(masterDB)
	methodValue := masterDBValue.MethodByName(methodName)
	if methodValue.IsValid() {
		methodValue.Call(in)
	}
	// slaveDBs call
	for _, slaveDB := range slaveDBs {
		slaveDBValue := reflect.ValueOf(slaveDB)
		methodValue := slaveDBValue.MethodByName(methodName)
		if methodValue.IsValid() {
			methodValue.Call(in)
		}
	}
}

func (pool *Pool) incrSlaveIdx() {
	pool.slaveIdx++
	if pool.slaveIdx >= len(pool.slaveDBs) {
		pool.slaveIdx = 0
	}
}

// Config establish a DB connection
func Config(dialect string, dbConfig DBConfig) {
	if dialect == "mysql" {
		db, err := sql.Open("mysql", dbConfig.Master)
		checkError(err)

		pool.masterDB = &huhDB{*db}

		for _, slaveDBAddr := range dbConfig.Slaves {
			db, err := sql.Open("mysql", slaveDBAddr)
			checkError(err)

			pool.slaveDBs = append(pool.slaveDBs, &huhDB{*db})
		}
		return
	}
	panic(ErrDialectNotSupported)
}

// Close DB connection
func Close() {
	callForAllDBs("Close")
}
