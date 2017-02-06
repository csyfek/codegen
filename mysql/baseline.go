package mysql


func Baseline() string{
return `
package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jackmanlabs/errors"
	"sync"
)

var (
	_db    *sql.DB
	_mutex sync.Mutex
)

func db() (*sql.DB, error) {
	_mutex.Lock()
	defer _mutex.Unlock()

	if _db != nil {
		return _db, nil
	}

	connString := "username:password@tcp(host:port)/database?parseTime=true"

	var err error
	_db, err = sql.Open("mysql", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	_db.SetMaxIdleConns(0)   // There are issues with MySQL/MariaDB and connection maintenance.
	_db.SetMaxOpenConns(100) // A nice, round number.

	return _db, nil
}

func tx() (*sql.Tx, error) {
	db, err := db()
	if err != nil {
		return nil, errors.Stack(err)
	}

	tx, err := db.Begin()
	if err != nil {
		return nil, errors.Stack(err)
	}

	return tx, nil
}
`
}