package mysql

var templateDelete string = `
var psDelete{{.model}} *sql.Stmt

func (this *DataSource) Delete{{.model}}(id string) error {

var err error

	if psDelete{{.model}} == nil{
		// language=MySQL
		q := {{template "templateDeleteSql" .}}

		psDelete{{.model}}, err = this.Prepare(q)
		if err != nil {
			return errs.Stack(err)
		}
	}
	args := []interface{}{id}

	_, err = psDelete{{.model}}.Exec(args...)
	if err != nil {
		return errs.Stack(err)
	}

	return nil
}
`

var templateDeleteTx string = `
func  (this *DataSource) Delete{{.model}}Tx(tx *sql.Tx, id string) error {
	// language=MySQL
	q := {{template "templateDeleteSql" .}}

	args := []interface{}{id}
	_, err := tx.Exec(q, args...)

	if err != nil {
		return errs.Stack(err)
	}

	return nil
}`

var templateDeleteSql string = "`" + `
DELETE FROM {{.table}}
{{if .members}}WHERE {{.table}}.{{with index .members 0}}{{.SqlName}}{{end}} = ?{{end}};` + "`"
