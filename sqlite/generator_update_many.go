package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func (this *generator) UpdateMany(pkgName string, def *codegen.Type) (string, error) {

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

func (this *generator) UpdateManyTx(pkgName string, def *codegen.Type) (string, error) {

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

var templateUpdateMany string = `
var psInsert{{.model}} *sql.Stmt

	func  (this *SqliteDataSource) Insert{{.model}}(x *{{.modelPackageName}}.{{.model}}) error {
	
var err error

if psInsert{{.model}} == nil{
	q := {{template "templateInsertSql" .}}

psInsert{{.model}}, err := this.Prepare(q)
		if err != nil {
			return errors.Stack(err)
		}
}

{{range .members}}
	{{if .SqlType}}{{else}}var x_{{.GoName}}[]byte{{end}}
{{end}}


{{range .members}}
	{{if .SqlType}}{{else}}
	x_{{.GoName}}, err = json.Marshal(x.{{.GoName}})
	if err != nil {
		return errors.Stack(err)
	}
{{end}}
{{end}}

	}
}

	args := []interface{}{
{{range .members}}
	{{if .SqlType}}
		&x.{{.GoName}},
	{{else}}
		&x_{{.GoName}},
{{end}}
{{end}}

	}

	_, err = psInsert{{.model}}.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	// nil is returned if no data was present.
	return nil

}
`

var templateUpdateManyTx string = `

	func  (this *SqliteDataSource) Insert{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) error {


	q := {{template "templateInsertSql" .}}


{{range .members}}
	{{if .SqlType}}{{else}}var x_{{.GoName}}[]byte{{end}}
{{end}}


{{range .members}}
	{{if .SqlType}}{{else}}
	x_{{.GoName}}, err = json.Marshal(x.{{.GoName}})
	if err != nil {
		return errors.Stack(err)
	}
{{end}}
{{end}}

	}
}

	args := []interface{}{
{{range .members}}
	{{if .SqlType}}
		&x.{{.GoName}},
	{{else}}
		&x_{{.GoName}},
{{end}}
{{end}}

	}

	_, err = tx.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	// nil is returned if no data was present.
	return nil

}
`

var templateSqlUpdateMany string = "`" + `

	INSERT INTO {{.table}} (
{{range $i, $member := .members}}
{{$member}}{{if last $i $}}{{else}},{{end}}
{{end}}
	) VALUES (
{{range $i, $member := .members}}
${{inc $i}}{{if last $i $}});{{else}},{{end}}
{{end}}
` + "`"
