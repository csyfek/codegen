package sqlite

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"log"
	"strings"
)

func (this *generator) Schema(parents []codegen.Parent) (string, error) {
	var (
		s string
	)

	for _, parent := range parents {

		var foreignKeys map[string]string = make(map[string]string)

		for i, member := range parent.Children {

			// types
			t, _ := sqlType(member.GoType)
			parent.Children[i].SqlType = t
			if (member.SqlName == "id" || strings.HasSuffix(member.SqlName, "_id")) && t == "TEXT" {
				parent.Children[i].SqlType = "CHAR(36)"
			}

			// constraints
			if i == 0 && member.SqlName == "id" {
				parent.Children[i].SqlConstraint = "PRIMARY KEY"
			} else {
				parent.Children[i].SqlConstraint = "NOT NULL"
			}

			// foreign keys
			if strings.HasSuffix(member.SqlName, "_id") {
				table := strings.TrimSuffix(member.SqlName, "_id")
				foreignKeys[member.SqlName] = table
			}
		}

		data := map[string]interface{}{
			"members":     parent.Children,
			"model":       parent.Name,
			"parents":     codegen.Plural(parent.Name),
			"table":       codegen.Plural(snaker.CamelToSnake(parent.Name)),
			"type":        parent.Name,
			"foreignKeys": foreignKeys,
		}

		subPatterns := map[string]string{}

		s_, err := codegen.Render(templateSchema, subPatterns, data)
		if err != nil {
			return "", errors.Stack(err)
		}

		s += s_
		log.Print("model:", parent.Name)
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
