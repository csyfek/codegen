package controler

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/common"
)

func GetCollection(def *common.Type) string {

	model := def.Name
	b := bytes.NewBuffer(nil)
	models := plural(model)

	fmt.Fprintf(b, "func Get%s(filters map[string]interface{}) ([]types.%s, error) {\n", models, model)
	fmt.Fprintf(b, "z, err := data.Get%s()\n", models)
	fmt.Fprint(b, `
	if err != nil {
		return z, errors.Stack(err)
	}

	return z, nil
}
	`)

	return b.String()
}
