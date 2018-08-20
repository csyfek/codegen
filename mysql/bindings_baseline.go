package mysql

import "fmt"

func (this *generator) BindingsBaseline(pkgName string) string {
	return fmt.Sprintf(`
package %s

import (
	"database/sql"
	errs errs "github.com/jackmanlabs/errors"
	_ "github.com/go-sql-driver/mysql"
)

type DataSource struct {
	*sql.DB
}

func New() (*DataSource, error) {

	connString := "username:password@tcp(host:port)/database?parseTime=true"

	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, errs.Stack(err)
	}

	ds := &DataSource{
		DB: db,
	}

	return ds, nil
}

`, pkgName)
}
