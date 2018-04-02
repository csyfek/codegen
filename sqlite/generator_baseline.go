package sqlite

import "fmt"

func (this *generator) Baseline(pkgName string) string {
	return fmt.Sprintf(`
package %s

import (
	"database/sql"
	"github.com/jackmanlabs/errors"
	_ "github.com/mattn/go-sqlite3"
)

type SqliteDataSource struct {
	sql.DB
}

func New() (*SqliteDataSource, error) {
	connString := "file::memory:?mode=memory&cache=shared"

	var err error
	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	sqliteDb := &SqliteDataSource{
		DB:*db,
	}

	return sqliteDb, nil
}
`, pkgName)
}
