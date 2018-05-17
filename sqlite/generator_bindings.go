package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func (this *generator) Bindings(importPaths []string, bindingsPackageName string, modelPackageName string, def *codegen.Model) (string, error) {

	var (
		err error
	)

	for i, member := range def.Members {
		def.Members[i].SqlType, _ = sqlType(member.GoType)
	}

	data := map[string]interface{}{
		"importPaths":         importPaths,
		"members":             def.Members,
		"model":               def.Name,
		"models":              codegen.Plural(def.Name),
		"table":               snaker.CamelToSnake(def.Name),
		"type":                def.Name,
		"bindingsPackageName": bindingsPackageName,
		"modelPackageName":    modelPackageName,
	}

	subPatterns := map[string]string{
		"templateDelete":        templateDelete,
		"templateDeleteSql":     templateDeleteSql,
		"templateDeleteTx":      templateDeleteTx,
		"templateInsertOne":     templateInsertOne,
		"templateInsertOneTx":   templateInsertOneTx,
		"templateInsertSql":     templateInsertSql,
		"templateSelectMany":    templateSelectMany,
		"templateSelectManySql": templateSelectManySql,
		"templateSelectManyTx":  templateSelectManyTx,
		"templateSelectOne":     templateSelectOne,
		"templateSelectOneSql":  templateSelectOneSql,
		"templateSelectOneTx":   templateSelectOneTx,
		"templateUpdateOne":     templateUpdateOne,
		"templateUpdateOneSql":  templateUpdateOneSql,
		"templateUpdateOneTx":   templateUpdateOneTx,
	}

	s, err := codegen.Render(templateBindings, subPatterns, data)
	if err != nil {
		return "", errors.Stack(err)
	}

	return s, nil

}

var templateBindings string = `
package {{.bindingsPackageName}}

import (
	"database/sql"
	"github.com/jackmanlabs/errors"
	{{range .importPaths}}"{{.}}"{{end}}
)

//##############################################################################
// TABLE: {{.table}}
// TYPE:  {{.type}}
//##############################################################################

/*============================================================================*/

{{template "templateSelectOne" .}}

/*============================================================================*/

{{template "templateSelectOneTx" .}}

/*============================================================================*/

{{template "templateSelectMany" .}}

/*============================================================================*/

{{template "templateSelectManyTx" .}}

/*============================================================================*/

{{template "templateInsertOne" .}}

/*============================================================================*/

{{template "templateInsertOneTx" .}}

/*============================================================================*/

{{template "templateUpdateOne" .}}

/*============================================================================*/

{{template "templateUpdateOneTx" .}}

/*============================================================================*/

{{template "templateDelete" .}}

/*============================================================================*/
	`
