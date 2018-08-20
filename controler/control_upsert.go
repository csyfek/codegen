package controler

import (
	"bytes"
	"fmt"

	"github.com/jackmanlabs/codegen"
)

func Upsert(def *codegen.Model) string {

	model := def.Name
	handler := bytes.NewBuffer(nil)

	// TODO: Find a better way of determining the primary key.
	pkey := def.Members[0]

	fmt.Fprintf(handler, "func Upsert%s(id string, x *types.%s) (*types.%s, error) {\n", model, model, model)
	fmt.Fprintf(handler, "if id != x.%s {\n", pkey.GoName)
	fmt.Fprint(handler, `

		return nil, errors.New("ID in URL parameter and object are incongruent.")
	}

	if uid := uuid.Parse(id); uid == nil {
		return nil, errors.New("ID provided is not a valid UUID.")
	}

	tx, err := data.Tx()
	if err != nil {
		return nil, errors.Stack(err)
	}

	`)

	fmt.Fprintf(handler, "x_, err := data.Get%sTx(tx, id)\n", model)

	fmt.Fprint(handler, `
	if err != nil {
		tx.Rollback()
		return nil, errors.Stack(err)
	}

	if x_ == nil {
	`)

	fmt.Fprintf(handler, "err = data.Insert%sTx(tx, x)\n", model)

	fmt.Fprint(handler, `
		if err != nil {
			tx.Rollback()
			return nil, errors.Stack(err)
		}
	} else {
		err = data.InsertAIURoomTx(tx, x)
		if err != nil {
			tx.Rollback()
			return nil, errors.Stack(err)
		}
	}

	x_, err = data.GetAIURoomTx(tx, x.BuildingCode)
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
