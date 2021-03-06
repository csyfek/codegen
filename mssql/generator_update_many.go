package mssql

import (
	"bytes"
	"fmt"

	"github.com/jackmanlabs/codegen"
)

func (this *generator) UpdateMany(pkgName string, def *codegen.Model) string {

	var (
		b             = bytes.NewBuffer(nil)
		funcName      = fmt.Sprintf("UpdateOne%ss", def.Name)
		funcNameSlave = fmt.Sprintf("Update%sTx", def.Name)
	)

	fmt.Fprintf(b, "func %s(z []%s.%s) error {\n", funcName, pkgName, def.Name)
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

func (this *generator) UpdateManyTx(pkgName string, def *codegen.Model) string {

	var (
		b             = bytes.NewBuffer(nil)
		funcName      = fmt.Sprintf("UpdateOne%ssTx", def.Name)
		funcNameSlave = fmt.Sprintf("Update%sTx", def.Name)
	)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, z []%s.%s) error {\n", funcName, pkgName, def.Name)
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
