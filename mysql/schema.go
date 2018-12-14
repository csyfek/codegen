package mysql

import (
	"bytes"
	"github.com/jackmanlabs/codegen/util"
	"log"
	"strings"

	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func (this *generator) Schema(pkg *codegen.Package) ([]byte, error) {
	var (
		b *bytes.Buffer = bytes.NewBuffer(nil)
	)

	for _, def := range pkg.Models {

		var foreignKeys map[string]string = make(map[string]string)

		for i, member := range def.Members {

			// types
			def.Members[i].SqlType, _ = sqlType(member.GoType)
			if (member.SqlName == "id" || strings.HasSuffix(member.SqlName, "_id")) && member.GoType == "string" {
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
				table = util.Plural(table)
				foreignKeys[member.SqlName] = table
			}
		}

		data := map[string]interface{}{
			"members":     def.Members,
			"model":       def.Name,
			"models":      util.Plural(def.Name),
			"table":       util.Plural(snaker.CamelToSnake(def.Name)),
			"type":        def.Name,
			"foreignKeys": foreignKeys,
		}

		subPatterns := map[string]string{}

		s, err := util.Render(templateSchema, subPatterns, data)
		if err != nil {
			return nil, errors.Stack(err)
		}

		_, err = b.Write(s)
		if err != nil {
			return nil, errors.Stack(err)
		}

		log.Print("model:", def.Name)
	}

	return b.Bytes(), nil
}

var templateSchema string = `
DROP TABLE IF EXISTS {{.table}};
CREATE TABLE {{.table}} (
{{range $i, $member := .members}}	{{.SqlName}}	{{.SqlType}} {{.SqlConstraint}}{{if last $i $.members}}{{else}},
{{end}}{{end}}{{range $col, $table := .foreignKeys}},
FOREIGN KEY ({{$col}}) REFERENCES {{$table}} (id) ON DELETE CASCADE{{end}}
);
`
