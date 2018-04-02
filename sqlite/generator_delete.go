package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func (this *generator) Delete(def *codegen.Type) (string, error) {

	data := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	s, err := render(templateDelete, map[string]string{"templateDeleteSql": templateDeleteSql}, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil
}

func (this *generator) DeleteTx(def *codegen.Type) (string, error) {

	data := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	s, err := render(templateDeleteTx, map[string]string{"templateDeleteSql": templateDeleteSql}, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil
}

var templateDelete string = `
var psDelete{{.model}} *sql.Stmt

func (this *SqliteDataSource) Delete{{.model}}(id string) error {

var err error

	if psDelete{{.model}} == nil{
		q := {{template "templateDeleteSql" .}}

		psDelete{{.model}}, err = this.Prepare(q)
		if err != nil {
			return errors.Stack(err)
		}
	}
	args := []interface{}{id}

	_, err = psDelete{{.model}}.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

var templateDeleteTx string = `
func  (this *SqliteDataSource) Delete{{.model}}Tx(tx *sql.Tx, id string) error {
	q := {{template "templateDeleteSql" .}}

	args := []interface{}{id}
	_, err := tx.Exec(q, args...)")

	if err != nil {
		return errors.Stack(err)
	}

	return nil
}`

var templateDeleteSql string = "`" + `
DELETE FROM {{.table}}
{{if .members}}WHERE {{.table}}.{{with index .members 0}}{{.SqlName}}{{end}} = $1{{end}}
;` + "`"
