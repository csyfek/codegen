package codegen

import (
	"github.com/jackmanlabs/codegen/util"
	"github.com/jackmanlabs/errors"
)

func PackageInterface(
	classes []*Model,
	interfacePackageName string,
) ([]byte, error) {

	data := map[string]interface{}{
		"classes":              classes,
		"interfacePackageName": interfacePackageName,
	}

	t := `
package {{.interfacePackageName}}
 
type DataSource interface{
{{range .classes}}	{{.Name}}DataSource
{{end}}
}`

	s, err := util.Render(t, map[string]string{}, data)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return s, nil
}

func ModelInterface(
	importPaths []string,
	interfacePackageName string,
	modelPackageName string,
	def *Model,
) ([]byte, error) {
	data := map[string]interface{}{
		"model":                def.Name,
		"modelPackageName":     modelPackageName,
		"interfacePackageName": interfacePackageName,
		"importPaths":          importPaths,
		"models":               util.Plural(def.Name),
	}

	t := `
package {{.interfacePackageName}}

import (
	"database/sql"
	{{range .importPaths}}"{{.}}"{{end}}
)
 
type {{.model}}DataSource interface{
	Delete{{.model}}(id string) error
	Delete{{.model}}Tx(tx *sql.Tx, id string) error 
	Insert{{.model}}(x *{{.modelPackageName}}.{{.model}}) (*{{.modelPackageName}}.{{.model}}, error)
	Insert{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) (*{{.modelPackageName}}.{{.model}}, error)
	Select{{.models}}() ([]{{.modelPackageName}}.{{.model}}, error) 
	Select{{.models}}Tx(tx *sql.Tx)  ([]{{.modelPackageName}}.{{.model}}, error) 
	Select{{.model}}(id string) (*{{.modelPackageName}}.{{.model}}, error)
	Select{{.model}}Tx(tx *sql.Tx, id string)  (*{{.modelPackageName}}.{{.model}}, error)
	Update{{.model}}(x *{{.modelPackageName}}.{{.model}}) error 
	Update{{.model}}Tx(tx *sql.Tx, x *{{.modelPackageName}}.{{.model}}) error
}`

	s, err := util.Render(t, map[string]string{}, data)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return s, nil
}
