package sqlite

func Baseline() string {
	return `
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/jackmanlabs/errors"
	_ "github.com/mattn/go-sqlite3"
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

	connString := "file::memory:?mode=memory&cache=shared"

	var err error
	_db, err = sql.Open("sqlite", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	_db.SetMaxIdleConns(10)
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
