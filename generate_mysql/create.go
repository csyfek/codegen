package generate_mysql

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

	var firstField GoSqlDatum
	if len(members) > 0 {
		firstField = members[0]
	}

	fmt.Fprintf(b, "CREATE TABLE %s (\n", tableName)
	for idx, member := range members {
		if idx == 0 {
			if member.Type == "string" {
				member.SqlType = "CHAR(36)"
			}
			fmt.Fprintf(b, "\t%s %s PRIMARY KEY,\n", member.SqlName, member.SqlType)
		} else if idx == len(members)-1 {
			fmt.Fprintf(b, "\t%s %s NOT NULL\n", member.SqlName, member.SqlType)
		} else {
			fmt.Fprintf(b, "\t%s %s NOT NULL,\n", member.SqlName, member.SqlType)
		}
	}
	fmt.Fprintf(b, "\t-- FOREIGN KEY (%s) REFERENCES parent_table (id) ON DELETE CASCADE\n", firstField.SqlName)
	fmt.Fprintf(b, ")\n")
	fmt.Fprint(b, "\tENGINE = TokuDB\n")
	fmt.Fprint(b, "\t-- ENGINE = InnoDB\n")
	fmt.Fprint(b, "\t-- ROW_FORMAT = COMPRESSED; -- Requires Barracuda file format.\n")
	fmt.Fprint(b, "\tDEFAULT CHARSET = utf8;\n")

	return b.String()
}
