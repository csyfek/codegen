package generate_mysql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/serenize/snaker"
)

func Delete(def structfinder.StructDefinition) string {

	members := getGoSqlData(def.Members)

	b := bytes.NewBuffer(nil)
	b_sql := deleteSql(def, members)

	funcName := fmt.Sprintf("Delete%s", def.Name)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(id string) error {\n", funcName)
	fmt.Fprint(b, `
	db, err := db()
	if err != nil {
		return errors.Stack(err)
	}

`)
	fmt.Fprintf(b, "\tif %s == nil{\n", psName)
	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprintf(b, "\t\t%s, err = db.Prepare(q)", psName)
	fmt.Fprint(b, `
		if err != nil {
		return errors.Stack(err)
		}
`)
	fmt.Fprint(b, "	}\n\n") // end of prepared statement clause
	fmt.Fprint(b, "\targs := []interface{}{id}\n\n")
	fmt.Fprintf(b, "\t_, err = %s.Exec(args...)", psName)
	fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

`)

	fmt.Fprint(b, "\treturn nil\n")
	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func deleteSql(def structfinder.StructDefinition, members []GoSqlDatum) *bytes.Buffer {

	b := bytes.NewBuffer(nil)
	tableName := snaker.CamelToSnake(def.Name)

	fmt.Fprintf(b, "DELETE FROM %s\n", tableName)
	if len(members) > 0 {
		member := members[0]
		fmt.Fprintf(b, "\tWHERE %s.%s = ?;\n", tableName, member.SqlName)
	} else {
		fmt.Fprint(b, "\t-- Insert your filter criteria here.\n")
	}

	return b
}
