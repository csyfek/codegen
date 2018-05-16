package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func (this *generator) UpdateOne(pkgName string, def *codegen.Type) (string, error) {

	data := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	s, err := render(templateUpdateOne, map[string]string{"templateUpdateOneSql": templateUpdateOneSql}, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil
}

func (this *generator) UpdateOneTx(pkgName string, def *codegen.Type) (string, error) {

	data := map[string]interface{}{
		"model":   def.Name,
		"members": def.Members,
		"table":   snaker.CamelToSnake(def.Name),
	}

	s, err := render(templateUpdateOneTx, map[string]string{"templateUpdateOneSql": templateUpdateOneSql}, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil

}

var templateUpdateOne string = `
var psUpdate{{.model}} *sql.Stmt

func  (this *SqliteDataSource) Update{{.model}}(x *{{.modelPackageName}}.{{.model}}) error {

var err error

if psUpdate{{.model}} == nil{
	q := {{template "templateUpdateOneSql" .}}

psUpdate{{.model}}, err = this.Prepare(q)
		if err != nil {
			return errors.Stack(err)
		}
}

	args := []interface{}{
{{range .members}}&x.{{.GoName}},
{{end}}
{{range $i, $member := .members}}{{if eq $i 0}}&x.{{$member.GoName}}{{end}}{{end}},
	}

	_, err = psUpdate{{.model}}.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

var templateUpdateOneTx string = `

	func  (this *SqliteDataSource) Update{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) error {




var err error

	q := {{template "templateUpdateOneSql" .}}


	args := []interface{}{
{{range .members}}&x.{{.GoName}},
{{end}}
{{range $i, $member := .members}}{{if eq $i 0}}&x.{{$member.GoName}}{{end}}{{end}},
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

var templateUpdateOneSql string = "`" + `
UPDATE {{.table}}
SET
{{range $i, $member := .members}}{{$member.SqlName}} = ?{{if last $i $}}{{else}},{{end}}
{{end}}
WHERE {{range $i, $member := .members}}{{if eq $i 0}}{{$member.SqlName}}{{end}}{{end}} = ?;
` + "`"
