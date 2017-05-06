package mssql

import (
	"bytes"
	"fmt"
)

func (this *Generator)UpdateOne(pkgName, typeName, table string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := updateSql(table, columns)

	funcName := fmt.Sprintf("UpdateOne%s", typeName)
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
	if len(columns) > 0 {
		fmt.Fprintf(b, "\t\t&x.%s,\n", columns[0].ColumnName)
	}
	fmt.Fprint(b, "\t}\n\n")

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

func  (this *Generator)UpdateOneTx(pkgName, typeName, table string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := updateSql(table, columns)

	funcName := fmt.Sprintf("Update%sTx", typeName)

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
	if len(columns) > 0 {
		fmt.Fprintf(b, "\t\t&x.%s,\n", columns[0].ColumnName)
	}
	fmt.Fprint(b, "\t}\n\n")

	fmt.Fprint(b, "\t_, err = tx.Exec(q, args...)")
	fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

`)

	fmt.Fprint(b, "\treturn nil\n")
	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

// I have to leave out back ticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and column names.
func updateSql(table string, columns []Column) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	var firstField Column
	if len(columns) > 0 {
		firstField = columns[0]
	}

	fmt.Fprintf(b, "UPDATE %s\n", table)
	fmt.Fprint(b, "SET\n")
	for idx, column := range columns {
		if idx == len(columns)-1 {
			fmt.Fprintf(b, "\t%s.%s = ?\n", table, column.ColumnName)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s.%s = ?,\n", table, column.ColumnName)
		}
	}
	fmt.Fprintf(b, "WHERE %s.%s = ?;\n", table, firstField.ColumnName)

	return b
}
