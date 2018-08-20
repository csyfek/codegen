package mysql

var templateInsertManyTx string = `

func (this *DataSource) Insert{{.models}}Tx(tx *sql.Tx, z []{{.modelPackageName}}.{{.model}}) error {

	var err error

	for i := range z {
		err = this.Insert{{.model}}Tx(tx, &z[i])
		if err != nil {
			return errs.Stack(err)
		}
	}

	return nil
}
`
