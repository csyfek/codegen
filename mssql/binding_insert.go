package mssql

import (
	"bytes"
	"fmt"
)

func (this *Generator) Insert(pkgName, typeName, tableName string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := insertSql(tableName, columns)

	funcName := fmt.Sprintf("Insert%s", typeName)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(x *%s.%s) error {\n", funcName, pkgName, typeName)
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

	fmt.Fprint(b, "\n")

	fmt.Fprint(b, "\targs := []interface{}{\n")
	for _, column := range columns {
		fmt.Fprintf(b, "\t\t&x.%s,\n", column.ColumnName)
	}
	fmt.Fprint(b, "\t}\n\n")

	fmt.Fprintf(b, "\t_, err = %s.Exec(args...)", psName)
	fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

`)

	fmt.Fprint(b, "\t// nil is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

func (this *Generator) InsertTx(pkgName, typeName, table string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := insertSql(table, columns)

	funcName := fmt.Sprintf("Insert%sTx", typeName)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, x *%s.%s) error {\n", funcName, pkgName, typeName)
	fmt.Fprint(b, "var err error\n")

	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprint(b, "\n")

	fmt.Fprint(b, "\targs := []interface{}{\n")
	for _, column := range columns {
		fmt.Fprintf(b, "\t\t&x.%s,\n", column.ColumnName)
	}
	fmt.Fprint(b, "\t}\n\n")

	fmt.Fprint(b, "\t_, err = tx.Exec(q, args...)")
	fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

`)

	fmt.Fprint(b, "\t// nil is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func insertSql(table string, columns []Column) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	fmt.Fprintf(b, "INSERT INTO %s (\n", table)
	for idx, column := range columns {
		if idx == len(columns)-1 {
			fmt.Fprintf(b, "\t%s\n", column.ColumnName)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s,\n", column.ColumnName)
		}
	}
	fmt.Fprint(b, ") VALUES (")
	for idx := range columns {
		if idx == len(columns)-1 {
			fmt.Fprint(b, "?);\n")
		} else {
			fmt.Fprint(b, "?, ")
		}
	}

	return b
}
