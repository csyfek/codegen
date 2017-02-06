package pg

// TODO: Supprt https://github.com/jackc/pgx

func Baseline() string {
	return `
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/jackmanlabs/errors"
	_ "github.com/lib/pq"
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

	connString := "postgres://username:password@localhost/database?sslmode=verify-full"

	var err error
	_db, err = sql.Open("postgres", connString)
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
