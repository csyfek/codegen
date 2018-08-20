package rester

import (
	"bytes"
	"fmt"

	"github.com/jackmanlabs/codegen"
)

func GetOne(def *codegen.Model) (string, string) {

	resourceName := resource(def.Name)
	model := def.Name

	register := fmt.Sprintf(`r.Path("/%s/{id}").Methods("GET").Handler(ErrFilter(handleGet%s))`, resourceName, model)

	handler := bytes.NewBuffer(nil)

	fmt.Fprintf(handler, "func handleGet%s(w http.ResponseWriter, r *http.Request) error {\n", model)
	fmt.Fprintln(handler, `var id string = mux.Vars(r)["id"]`)
	fmt.Fprintf(handler, "x, err := control.Get%s(id)\n", model)
	fmt.Fprint(handler, `
	if err != nil {
		return errors.Stack(err)
	}else if x == nil {
		w.WriteHeader(404)
		return nil
	}

	err = serialize(w,r,x)
	if err != nil{
		return errors.Stack(err)
	}

	return nil
}
	`)

	return register, handler.String()
}
