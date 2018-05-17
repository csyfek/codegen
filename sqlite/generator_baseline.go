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

type DataSource struct {
	*sql.DB
}

func New() (*DataSource, error) {
	connString := "file::memory:?mode=memory&cache=shared"

	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	ds := &DataSource{
		DB:db,
	}

	return ds, nil
}
`, pkgName)
}
