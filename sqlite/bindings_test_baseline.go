package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
)

func (this *generator) BindingsBaselineTests(importPaths []string, bindingsPackageName string, modelPackageName string) (string, error) {

	var (
		err error
	)


	data := map[string]interface{}{
		"importPaths":         importPaths,
		"bindingsPackageName": bindingsPackageName,
		"modelPackageName":    modelPackageName,
	}

	subPatterns := map[string]string{
	}

	s, err := codegen.Render(templateTestBaseline, subPatterns, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil
}

var templateTestBaseline string = `
package {{.bindingsPackageName}}_test

import (
{{range .importPaths}}	"{{.}}"
{{end}}
	"database/sql"
	"github.com/jackmanlabs/errors"
	_ "github.com/mattn/go-sqlite3"
)

// This is redundant, I know, but I don't think it's wise to count on the New()
// function of the bindings package to function to behave the same as it was
// written by this generator. So... we're going to assume that, at minimum, the
// type remains a wrapper of a sql.DB.
func New() (*{{.bindingsPackageName}}.DataSource, error) {
	connString := "file::memory:?mode=memory&cache=shared"

	db, err := sql.Open("sqlite3", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	ds := &{{.bindingsPackageName}}.DataSource{
		DB:db,
	}

	return ds, nil
}
`
