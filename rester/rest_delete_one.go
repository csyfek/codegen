package rester

import (
	"bytes"
	"fmt"

	"github.com/jackmanlabs/codegen"
)

func DeleteOne(def *codegen.Model) (string, string) {

	resourceName := resource(def.Name)
	model := def.Name

	register := fmt.Sprintf(`r.Path("/%s/{id}").Methods("DELETE").Handler(ErrFilter(handleDelete%s))`, resourceName, model)

	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func handleDelete%s(w http.ResponseWriter, r *http.Request) error {\n", model)
	fmt.Fprintln(handler, `var id string = mux.Vars(r)["id"]`)
	fmt.Fprintf(handler, "err := control.Delete%s(id)\n", model)
	fmt.Fprint(handler, `
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}`)

	return register, handler.String()
}
