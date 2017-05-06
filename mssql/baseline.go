package mssql


func (this *Generator) Baseline() string {
	return `
package data

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
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

	connString := "server=host;database=database;user id=username;password=password"

	var err error
	_db, err = sql.Open("mssql", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

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

// This is available assuming that you, like me, want to keep advanced DB
// operations in another package.
func Tx() (*sql.Tx, error) {
	return tx()
}
`
}

