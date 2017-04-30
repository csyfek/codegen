package mysql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/extractor"
)

func UpdatePlural(pkgName string, def *extractor.StructDefinition) string {

	var (
		b             = bytes.NewBuffer(nil)
		funcName      = fmt.Sprintf("Update%ss", def.Name)
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
		return nil, errors.Stack(err)
	}

	return nil
}`)

	return b.String()
}


func UpdatePluralTx(pkgName string, def *extractor.StructDefinition) string {

	var (
		b             = bytes.NewBuffer(nil)
		funcName      = fmt.Sprintf("Update%ssTx", def.Name)
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