package sqlite

import (
	"bytes"
	"github.com/jackmanlabs/errors"
	"text/template"
	"reflect"
)

func render(rootPattern string, subPatterns map[string]string, data interface{}) (string, error) {

	var fns = template.FuncMap{
		"last": func(x int, a interface{}) bool {
			return x == reflect.ValueOf(a).Len()-1
		},
		"inc": func(x int) int { return x + 1 },
	}

	t, err := template.New("").Funcs(fns).Parse(rootPattern)
	if err != nil {
		return "", errors.Stack(err)
	}

	for name, pattern := range subPatterns {
		_, err := t.New(name).Parse(pattern)
		if err != nil {
			return "", errors.Stack(err)
		}
	}

	b := bytes.NewBuffer(nil)

	err = t.Execute(b, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return b.String(), nil

}
