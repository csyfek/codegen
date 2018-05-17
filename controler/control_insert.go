package controler

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func Insert(def *codegen.Model) string {

	// TODO: Find a better way of determining the primary key.
	pkey := def.Members[0]

	model := def.Name
	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func Insert%s(x *types.%s) (*types.%s, error) {\n", model, model, model)
	//fmt.Fprintf(handler, "x.%s= uuid.New()\n", index.ColumnName)

	fmt.Fprint(handler, `
	tx, err := data.Tx()
	if err != nil {
		return nil, errors.Stack(err)
	}
	`)

	fmt.Fprintf(handler, "err = data.Insert%sTx(tx, x)\n", model)
	fmt.Fprint(handler, `
	if err != nil {
		tx.Rollback()
		return nil, errors.Stack(err)
	}
`)

	fmt.Fprintf(handler, "x_, err := data.Get%sTx(tx, x.%s)\n", model, pkey.GoName)
	fmt.Fprint(handler, `
	if err != nil {
		tx.Rollback()
		return nil, errors.Stack(err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, errors.Stack(err)
	}

	return x_, nil
}
	`)

	return handler.String()
}
