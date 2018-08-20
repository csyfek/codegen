package mysql

var templateSelectMany string = `
var psSelect{{.models}} *sql.Stmt

func (this *DataSource)  Select{{.models}}() ([]{{.modelPackageName}}.{{.model}}, error) {

	var err error

	if psSelect{{.models}} == nil{
		// language=MySQL
		q := {{template "templateSelectManySql" .}}
	
		psSelect{{.models}}, err = this.Prepare(q)
		if err != nil {
			return nil, errs.Stack(err)
		}
	}

	args := []interface{}{}

	rows, err := psSelect{{.models}}.Query( args...)
	if err != nil {
		return nil, errs.Stack(err)
	}
	defer rows.Close()

	var z []{{.modelPackageName}}.{{.model}} = make([]{{.modelPackageName}}.{{.model}},0)
	for rows.Next() {
		var x {{.modelPackageName}}.{{.model}}
		dest := []interface{}{
{{range .members}}			&x.{{.GoName}},
{{end}}		}

		err = rows.Scan(dest...)
		if err != nil {
			return z, errs.Stack(err)
		}

		z = append(z, x)
	}

	return z, nil
}
`

var templateSelectManyTx string = `

func  (this *DataSource) Select{{.models}}Tx(tx *sql.Tx)  ([]{{.modelPackageName}}.{{.model}}, error) {

	// language=MySQL
	q := {{template "templateSelectManyTxSql" .}}

	args := []interface{}{}

	rows, err := tx.Query(q, args...)
	if err != nil {
		return nil, errs.Stack(err)
	}
	defer rows.Close()

	var z []{{.modelPackageName}}.{{.model}} = make([]{{.modelPackageName}}.{{.model}},0)
	for rows.Next() {
		var x {{.modelPackageName}}.{{.model}}
		dest := []interface{}{
{{range .members}}			&x.{{.GoName}},
{{end}}		}

		err = rows.Scan(dest...)
		if err != nil {
			return z, errs.Stack(err)
		}

		z = append(z, x)
	}

	return z, nil
}
`

var templateSelectManySql string = "`" + `
SELECT
{{range $i, $member := .members}}	{{$member.SqlName}}{{if last $i $.members}}{{else}},
{{end}}{{end}}
FROM {{.table}};
` + "`"

var templateSelectManyTxSql string = "`" + `
SELECT
{{range $i, $member := .members}}	{{$member.SqlName}}{{if last $i $.members}}{{else}},
{{end}}{{end}}
FROM {{.table}}
FOR UPDATE;
` + "`"
