package rester

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func PutOne(def *codegen.Type) (string, string) {

	resourceName := resource(def.Name)
	model := def.Name

	register := fmt.Sprintf(`r.Path("/%s/{id}").Methods("PUT").Handler(ErrFilter(handlePut%s))`, resourceName, model)

	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func handlePut%s(w http.ResponseWriter, r *http.Request) error {\n", model)
	fmt.Fprintln(handler, `var id string = mux.Vars(r)["id"]`)
	fmt.Fprintf(handler, "var x *types.%s = new (types.%s)\n", model, model)
	fmt.Fprint(handler, `
	var err error
	err = deserialize(r,x)
	if err != nil{
		return errors.Stack(err)
	}

	`)

	fmt.Fprintf(handler, "x_, err := control.Upsert%s(id, x)\n", model)
	fmt.Fprint(handler, `
	if err != nil {
		return errors.Stack(err)
	}

	err = serialize(w,r,x_)
	if err != nil{
		return errors.Stack(err)
	}

	return nil
}`)

	return register, handler.String()
}
