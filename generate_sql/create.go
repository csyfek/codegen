package generate_sql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/serenize/snaker"
)

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func Create(def structfinder.StructDefinition) string {
	members := getGoSqlData(def.Members)

	b := bytes.NewBuffer(nil)
	tableName := snaker.CamelToSnake(def.Name)

	var firstField structfinder.StructMemberDefinition
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprintf(b, "CREATE TABLE %s (\n", tableName)
	for _, member := range members {
		fmt.Fprintf(b, "\t%s %s NOT NULL,\n", member.SqlName, member.SqlType)
	}
	fmt.Fprintf(b, "\tPRIMARY KEY (%s)\n", snaker.CamelToSnake(firstField.Name))
	fmt.Fprintf(b, ")\n")
	fmt.Fprint(b, "\tENGINE = InnoDB\n")
	fmt.Fprint(b, "\tDEFAULT CHARSET = utf8\n")
	fmt.Fprint(b, "\tROW_FORMAT = COMPRESSED;\n")

	return b.String()
}
