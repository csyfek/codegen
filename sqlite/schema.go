package sqlite

import (
	"log"
	"strings"

	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func (this *generator) Schema(pkg *codegen.Package) (string, error) {
	var (
		s string
	)

	for _, def := range pkg.Models {

		var foreignKeys map[string]string = make(map[string]string)

		for i, member := range def.Members {

			// types
			t, _ := sqlType(member.GoType)
			def.Members[i].SqlType = t
			if (member.SqlName == "id" || strings.HasSuffix(member.SqlName, "_id")) && t == "TEXT" {
				def.Members[i].SqlType = "CHAR(36)"
			}

			// constraints
			if i == 0 && member.SqlName == "id" {
				def.Members[i].SqlConstraint = "PRIMARY KEY"
			} else {
				def.Members[i].SqlConstraint = "NOT NULL"
			}

			// foreign keys
			if strings.HasSuffix(member.SqlName, "_id") {
				table := strings.TrimSuffix(member.SqlName, "_id")
				foreignKeys[member.SqlName] = table
			}
		}

		data := map[string]interface{}{
			"members":     def.Members,
			"model":       def.Name,
			"models":      codegen.Plural(def.Name),
			"table":       snaker.CamelToSnake(def.Name),
			"type":        def.Name,
			"foreignKeys": foreignKeys,
		}

		subPatterns := map[string]string{}

		s_, err := codegen.Render(templateSchema, subPatterns, data)
		if err != nil {
			return "", errors.Stack(err)
		}

		s += s_
		log.Print("model:", def.Name)
	}

	return s, nil
}

var templateSchema string = `
DROP TABLE IF EXISTS {{.table}};
CREATE TABLE {{.table}} (
{{range $i, $member := .members}}	{{.SqlName}}	{{.SqlType}} {{.SqlConstraint}}{{if last $i $.members}}{{else}},
{{end}}{{end}}{{range $col, $table := .foreignKeys}},
FOREIGN KEY ({{$col}}) REFERENCES {{$table}} (id) ON DELETE CASCADE{{end}}
);
`
