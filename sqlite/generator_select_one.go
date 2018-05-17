package sqlite

var templateSelectOne string = `
var psSelect{{.model}} *sql.Stmt

func (this *DataSource)  Select{{.model}}(id string) (*{{.modelPackageName}}.{{.model}}, error) {

	var err error

	if psSelect{{.model}} == nil{
		q := {{template "templateSelectOneSql" .}}
	
		psSelect{{.model}}, err = this.Prepare(q)
		if err != nil {
			return nil, errors.Stack(err)
		}
	}

	args := []interface{}{id}

	rows, err := psSelect{{.model}}.Query(args...)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

	var x *{{.modelPackageName}}.{{.model}}
	if rows.Next() {
		x = new({{.modelPackageName}}.{{.model}})
		dest := []interface{}{
{{range .members}}			&x.{{.GoName}},
{{end}}		}

		err = rows.Scan(dest...)
		if err != nil {
			return x, errors.Stack(err)
		}
	}

	return x, nil
}
`

var templateSelectOneTx string = `

func  (this *DataSource) Select{{.model}}Tx(tx *sql.Tx, id string)  (*{{.modelPackageName}}.{{.model}}, error) {

	q := {{template "templateSelectOneSql" .}}

	args := []interface{}{id}

	rows, err := tx.Query(q, args...)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

	var x *{{.modelPackageName}}.{{.model}}
	if rows.Next() {
		x = new({{.modelPackageName}}.{{.model}})
		dest := []interface{}{
{{range .members}}			&x.{{.GoName}},
{{end}}		}

		err = rows.Scan(dest...)
		if err != nil {
			return x, errors.Stack(err)
		}
	}

	return x, nil
}
`

var templateSelectOneSql string = "`" + `
SELECT
{{range $i, $member := .members}}	{{$member.SqlName}}{{if last $i $.members}}{{else}},{{end}}
{{end}}FROM {{.table}}
WHERE {{range $i, $member := .members}}{{if eq $i 0}}{{$member.SqlName}}{{end}}{{end}} = ?
LIMIT 1;
` + "`"
