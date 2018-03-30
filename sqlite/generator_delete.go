package sqlite

import (
	"bytes"
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"text/template"
)

func (this *generator) Delete(def *codegen.Type) (string, error) {

	values := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	pattern := `
var ps_Delete{{.model}} *sql.Stmt

func Delete{{.model}}(id string) error {

	db, err := db()
	if err != nil {
		return errors.Stack(err)
	}

	if ps_Delete{{.model}} == nil{
		q := {{template "sql" .}}

		ps_Delete{{.model}}, err = db.Prepare(q)
		if err != nil {
			return errors.Stack(err)
		}
	}
	args := []interface{}{id}

	_, err = ps_Delete{{.model}}.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

	b := bytes.NewBuffer(nil)

	t, err := template.New("go").Parse(pattern)
	if err != nil {
		return "", errors.Stack(err)
	}

	_, err = t.New("sql").Parse(templateSqlDelete)
	if err != nil {
		return "", errors.Stack(err)
	}

	err = t.Execute(b, values)
	if err != nil {
		return "", errors.Stack(err)
	}

	return b.String(), nil
}

func (this *generator) DeleteTx(def *codegen.Type) (string, error) {

	values := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	pattern := `
func Delete{{}}Tx(tx *sql.Tx, id string) error {
	q := {{template "sql" .}}

	args := []interface{}{id}
	_, err := tx.Exec(q, args...)")

	if err != nil {
		return errors.Stack(err)
	}

return nil
}`

	b := bytes.NewBuffer(nil)

	t, err := template.New("go").Parse(pattern)
	if err != nil {
		return "", errors.Stack(err)
	}

	_, err = t.New("sql").Parse(templateSqlDelete)
	if err != nil {
		return "", errors.Stack(err)
	}

	err = t.Execute(b, values)
	if err != nil {
		return "", errors.Stack(err)
	}

	return b.String(), nil

}

var templateSqlDelete string = "`" + `
DELETE FROM {{.table}}
{{if .members}}WHERE {{.table}}.{{with index .members 0}}{{.SqlName}}{{end}} = $1{{end}}
;` + "`"
