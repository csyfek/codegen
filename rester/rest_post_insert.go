package rester

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func PostCollection(def *codegen.Type) (string, string) {

	resourceName := resource(def.Name)
	model := def.Name

	register := fmt.Sprintf(`r.Path("/%s").Methods("POST").Handler(ErrFilter(handlePost%s))`, resourceName, model)

	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func handlePost%s(w http.ResponseWriter, r *http.Request) error {\n", model)
	fmt.Fprintf(handler, "var x *types.%s = new (types.%s)\n", model, model)
	fmt.Fprint(handler, `
	var err error
	err = deserialize(r,x)
	if err != nil{
		return errors.Stack(err)
	}

	`)

	fmt.Fprintf(handler, "x_, err := control.Insert%s(x)\n", model)
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
