package controler

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func Delete(def *codegen.Parent) string {

	model := def.Name
	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func Delete%s(id string) error {\n", model)
	fmt.Fprintf(handler, "err := data.Delete%s(id)\n", model)
	fmt.Fprint(handler, `
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}`)

	return handler.String()
}
