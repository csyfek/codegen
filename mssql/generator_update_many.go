package mssql

import (
	"bytes"
	"fmt"
)

func (this *generator) UpdateMany(pkgName, typeName, table string, columns []Column) string {

	var (
		b             = bytes.NewBuffer(nil)
		funcName      = fmt.Sprintf("UpdateOne%ss", typeName)
		funcNameSlave = fmt.Sprintf("Update%sTx", typeName)
	)

	fmt.Fprintf(b, "func %s(z []%s.%s) error {\n", funcName, pkgName, typeName)
	fmt.Fprint(b, `

	tx, err := tx()
	if err != nil{
		return errors.Stack(err)
	}

	for _, x := range z {
`)
	fmt.Fprintf(b, "err := %s(tx, &x)", funcNameSlave)
	fmt.Fprint(b, `
		if err != nil {
			return errors.Stack(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}`)

	return b.String()
}

func (this *generator) UpdateManyTx(pkgName, typeName, table string, columns []Column) string {

	var (
		b             = bytes.NewBuffer(nil)
		funcName      = fmt.Sprintf("UpdateOne%ssTx", typeName)
		funcNameSlave = fmt.Sprintf("Update%sTx", typeName)
	)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, z []%s.%s) error {\n", funcName, pkgName, typeName)
	fmt.Fprint(b, `

	for _, x := range z {
`)
	fmt.Fprintf(b, "err := %s(tx, &x)", funcNameSlave)
	fmt.Fprint(b, `
		if err != nil {
			return errors.Stack(err)
		}
	}

	return nil
}`)

	return b.String()
}
