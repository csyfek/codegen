package mysql



var templateInsertOne string = `
var psInsert{{.model}} *sql.Stmt

func  (this *DataSource) Insert{{.model}}(x *{{.modelPackageName}}.{{.model}}) error {

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

	func  (this *DataSource) Insert{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) error {

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
{{range $i, $member := .members}}	{{$member.SqlName}}{{if last $i $.members}}{{else}},{{end}}
{{end}}) VALUES (
{{range $i, $member := .members}}	?{{if last $i $.members}}{{else}},{{end}}
{{end}});
` + "`"
