package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/serenize/snaker"
	"github.com/jackmanlabs/errors"
)

func (this *generator) InsertOne(pkgName string, def *codegen.Type) (string, error) {

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

func (this *generator) InsertOneTx(pkgName string, def *codegen.Type) (string, error) {

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

var templateInsertOne string = `
var psInsert{{.model}} *sql.Stmt

func  (this *SqliteDataSource) Insert{{.model}}(x *{{.modelPackageName}}.{{.model}}) error {

	var err error

	if psInsert{{.model}} == nil{
		q := {{template "templateInsertSql" .}}

		psInsert{{.model}}, err = this.Prepare(q)
		if err != nil {
			return errors.Stack(err)
		}
	}

	args := []interface{}{
{{range .members}}&x.{{.GoName}},
{{end}}
	}

	_, err = psInsert{{.model}}.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

var templateInsertOneTx string = `

	func  (this *SqliteDataSource) Insert{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) error {

	var err error

	q := {{template "templateInsertSql" .}}

	args := []interface{}{
{{range .members}}&x.{{.GoName}},
{{end}}
	}

	_, err = tx.Exec(q, args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

var templateInsertSql string = "`" + `
INSERT INTO {{.table}} (
{{range $i, $member := .members}}{{$member.SqlName}}{{if last $i $}}{{else}},{{end}}
{{end}}) VALUES (
{{range $i, $member := .members}}${{inc $i}}{{if last $i $}});{{else}},{{end}}
{{end}}
` + "`"
