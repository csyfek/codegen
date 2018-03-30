package controler

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func GetOne(def *codegen.Type) string {

	model := def.Name
	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func Get%s(id string) (*types.%s, error) {\n", model, model)
	fmt.Fprintf(handler, "x_, err := data.Get%s(id)\n", model)
	fmt.Fprint(handler, `
	if err != nil {
		return x_, errors.Stack(err)
	}

	return x_, nil
}
	`)

	return handler.String()
}
