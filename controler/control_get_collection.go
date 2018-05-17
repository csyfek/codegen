package controler

import (
	"bytes"
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"text/template"
)

func GetCollection(def *codegen.Model) (string, error) {

	model := def.Name
	b := bytes.NewBuffer(nil)

	values := map[string]string{
		"singular": model,
		"plural":   codegen.Plural(model),
	}

	pattern := `
func Get{{.plural}}(filters map[string]interface{}) ([]types.{{.singular}}, error) {
	tz, err := data.Get{{.plural}}(filter filters.{{.singular}})
	if err != nil {
			return z, errors.Stack(err)
	}

	return z, nil
}
`

	tmpl, err := template.ParseGlob(pattern)
	if err != nil {
		return "", errors.Stack(err)
	}

	err = tmpl.Execute(b, values)
	if err != nil {
		return "", errors.Stack(err)
	}

	return b.String(), nil
}
