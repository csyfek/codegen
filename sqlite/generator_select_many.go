package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/serenize/snaker"
	"github.com/jackmanlabs/errors"
)

func (this *generator) SelectMany(pkgName string, def *codegen.Type) (string, error) {

	data := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	s, err := render(templateInsertOne, map[string]string{"templateInsertSql": templateInsertSql}, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil
}

func (this *generator) SelectManyTx(pkgName string, def *codegen.Type) (string, error) {

	data := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	s, err := render(templateInsertOneTx, map[string]string{"templateInsertSql": templateInsertSql}, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil

}

var templateSelectMany string = `
var psSelect{{.models}} *sql.Stmt

func (this *SqliteDataSource)  Select{{.models}}(id string) ([]{{.modelPackageName}}.{{.model}}, error) {

	var err error

	if psSelect{{.model}} == nil{
		q := {{template "templateSelectOneSql" .}}
	
		psSelect{{.model}}, err = this.Prepare(q)
		if err != nil {
			return nil, errors.Stack(err)
		}
	}

	args := []interface{}{id}

	rows, err := psSelect{{.models}}.Query( args...)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

	var z []{{.modelPackageName}}.{{.model}} = make([]{{.modelPackageName}}.{{.model}},0)
	for rows.Next() {
		var x {{.modelPackageName}}.{{.model}}
		dest := []interface{}{
{{range .members}}&x.{{.GoName}},
{{end}}
		}

		err = rows.Scan(dest...)
		if err != nil {
			return z, errors.Stack(err)
		}

		z = append(z, x)
	}

	return z, nil
}
`

var templateSelectManyTx string = `

func  (this *SqliteDataSource) Select{{.models}}Tx(tx *sql.Tx, id string)  ([]{{.modelPackageName}}.{{.model}}, error) {

	q := {{template "templateSelectOneSql" .}}

	args := []interface{}{id}

	rows, err := tx.Query(q, args...)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

	var z []{{.modelPackageName}}.{{.model}} = make([]{{.modelPackageName}}.{{.model}},0)
	for rows.Next() {
		var x {{.modelPackageName}}.{{.model}}
		dest := []interface{}{
{{range .members}}&x.{{.GoName}},
{{end}}
		}

		err = rows.Scan(dest...)
		if err != nil {
			return z, errors.Stack(err)
		}

		z = append(z, x)
	}

	return z, nil
}
`


var templateSelectManySql string = "`" + `
SELECT
{{range $i, $member := .members}}{{$member.SqlName}}{{if last $i $}}{{else}},
{{end}}{{end}}FROM {{.table}};
` + "`"
