package mysql

var templateUpdateOne string = `
var psUpdate{{.model}} *sql.Stmt

func  (this *DataSource) Update{{.model}}(x *{{.modelPackageName}}.{{.model}}) error {

var err error

if psUpdate{{.model}} == nil{
	q := {{template "templateUpdateOneSql" .}}

psUpdate{{.model}}, err = this.Prepare(q)
		if err != nil {
			return errors.Stack(err)
		}
}

	args := []interface{}{
{{range .members}}		&x.{{.GoName}},
{{end}}{{range $i, $member := .members}}{{if eq $i 0}}		&x.{{$member.GoName}}{{end}}{{end}},
	}

	_, err = psUpdate{{.model}}.Exec(args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`

var templateUpdateOneTx string = `

	func  (this *DataSource) Update{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) error {




var err error

	q := {{template "templateUpdateOneSql" .}}


	args := []interface{}{
{{range .members}}		&x.{{.GoName}},
{{end}}{{range $i, $member := .members}}{{if eq $i 0}}		&x.{{$member.GoName}}{{end}}{{end}},
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
{{range $i, $member := .members}}	{{$member.SqlName}} = ?{{if last $i $.members}}{{else}},{{end}}
{{end}}WHERE {{range $i, $member := .members}}{{if eq $i 0}}{{$member.SqlName}}{{end}}{{end}} = ?;
` + "`"
