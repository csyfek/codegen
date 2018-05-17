package sqlite

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
	"github.com/serenize/snaker"
)

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func (this *generator) Schema(pkg *codegen.Package) string {

	b := bytes.NewBuffer(nil)

	for _, def := range pkg.Models {
		fmt.Fprint(b, "\n\n/-- ----------------------------------------------------------------------------\n\n")
		b.WriteString(this.typeSchema(def))
	}

	return b.String()
}

func (this *generator) typeSchema(def *codegen.Model) string {
	b := bytes.NewBuffer(nil)
	tableName := snaker.CamelToSnake(def.Name)

	var firstField codegen.Member
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprintf(b, "DROP TABLE IF EXISTS %s;\n\n", tableName)
	fmt.Fprintf(b, "CREATE TABLE %s (\n", tableName)
	for idx, member := range def.Members {

		sqlType, _ := sqlType(member.GoType)

		if idx == 0 {
			if member.GoType == "string" {
				sqlType = "CHAR(36)"
			}
			fmt.Fprintf(b, "\t%s %s PRIMARY KEY,\n", member.SqlName, sqlType)
		} else if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s %s NOT NULL\n", member.SqlName, sqlType)
		} else {
			fmt.Fprintf(b, "\t%s %s NOT NULL,\n", member.SqlName, sqlType)
		}
	}
	fmt.Fprintf(b, "\t-- FOREIGN KEY (%s) REFERENCES parent_table (id) ON DELETE CASCADE\n", firstField.SqlName)
	fmt.Fprintf(b, ");\n")
	return b.String()
}
