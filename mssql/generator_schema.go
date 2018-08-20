package mssql

import (
	"bytes"
	"fmt"

	"github.com/jackmanlabs/codegen"
)

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func (this *generator) Schema(pkg *codegen.Package) string {

	b := bytes.NewBuffer(nil)

	for _, def := range pkg.Models {
		b.WriteString(this.typeDrop(def))
	}

	fmt.Fprint(b, "\n\n-- -----------------------------------------------------------------------------\n\n")

	for _, def := range pkg.Models {
		b.WriteString(this.typeSchema(def))
		fmt.Fprint(b, "\n\n-- -----------------------------------------------------------------------------\n\n")
	}

	return b.String()
}

func (this *generator) typeDrop(def *codegen.Model) string {
	return fmt.Sprintf("DROP TABLE IF EXISTS %s;\n", def.Table)
}

func (this *generator) typeSchema(def *codegen.Model) string {
	b := bytes.NewBuffer(nil)

	var firstField codegen.Member
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprintf(b, "CREATE TABLE %s (\n", def.Table)
	for idx, member := range def.Members {

		sqlType, _ := sqlType(member.GoType)

		if idx == 0 {
			//if member.GoType() == "string" {
			//	sqlType = "CHAR(36)"
			//}
			fmt.Fprintf(b, "\t%s %s PRIMARY KEY,\n", member.SqlName, sqlType)
		} else if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s %s NOT NULL\n", member.SqlName, sqlType)
		} else {
			fmt.Fprintf(b, "\t%s %s NOT NULL,\n", member.SqlName, sqlType)
		}
	}
	fmt.Fprintf(b, "\t-- FOREIGN KEY (%s) REFERENCES parent_table (id) ON DELETE CASCADE\n", firstField.SqlName)
	fmt.Fprint(b, ");\n")

	return b.String()
}
