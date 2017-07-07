package mssql

import (
	"bytes"
	"fmt"
)

func (this *generator) SelectOne(pkgName string, typeName string, table string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := selectOneSql(table, columns)

	funcName := fmt.Sprintf("Get%s", typeName)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(id string) (*%s.%s, error) {\n", funcName, pkgName, typeName)
	fmt.Fprint(b, `
	db, err := db()
	if err != nil {
		return nil, errors.Stack(err)
	}

`)
	fmt.Fprintf(b, "\tif %s == nil{\n", psName)
	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprintf(b, "\t\t%s, err = db.Prepare(q)", psName)
	fmt.Fprint(b, `
		if err != nil {
			return nil, errors.Stack(err)
		}
`)
	fmt.Fprint(b, "	}\n\n") // end of prepared statement clause
	fmt.Fprint(b, "\targs := []interface{}{id}\n\n")
	fmt.Fprintf(b, "\trows, err := %s.Query(args...)", psName)
	fmt.Fprint(b, `
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

`)

	fmt.Fprintf(b, "\tvar x *%s.%s\n", pkgName, typeName)
	fmt.Fprint(b, "\tif rows.Next() {\n")
	fmt.Fprintf(b, "\t\tx = new(%s.%s)\n", pkgName, typeName)

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, column := range columns {
		fmt.Fprintf(b, "\t\t\t&x.%s,\n", column.ColumnName)
	}

	fmt.Fprint(b, "\t\t}\n") // end of targets declaration.
	fmt.Fprint(b, `
		err = rows.Scan(targets...)
		if err != nil {
			return x, errors.Stack(err)
		}

`)

	fmt.Fprint(b, "\t}\n\n") // end of scan clause.
	fmt.Fprint(b, "\t// nil is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn x, nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

func (this *generator) SelectOneTx(pkgName, typeName, tableName string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := selectOneSqlTx(tableName, columns)

	funcName := fmt.Sprintf("Get%sTx", typeName)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, id string) (*%s.%s, error) {\n", funcName, pkgName, typeName)
	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprint(b, "\targs := []interface{}{id}\n\n")
	fmt.Fprint(b, "\trows, err := tx.Query(q, args...)")
	fmt.Fprint(b, `
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

`)

	fmt.Fprintf(b, "\tvar x *%s.%s\n", pkgName, typeName)
	fmt.Fprint(b, "\tif rows.Next() {\n")
	fmt.Fprintf(b, "\t\tx = new(%s.%s)\n", pkgName, typeName)

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, column := range columns {
		fmt.Fprintf(b, "\t\t\t&x.%s,\n", column.ColumnName)
	}

	fmt.Fprint(b, "\t\t}\n") // end of targets declaration.
	fmt.Fprint(b, `
		err = rows.Scan(targets...)
		if err != nil {
			return x, errors.Stack(err)
		}

`)

	fmt.Fprint(b, "\t}\n\n") // end of scan clause.
	fmt.Fprint(b, "\t// nil is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn x, nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

// I have to leave out back ticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and column names.
func selectOneSql(tableName string, columns []Column) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	var firstField Column
	if len(columns) > 0 {
		firstField = columns[0]
	}

	fmt.Fprint(b, "SELECT\n")
	for idx, column := range columns {
		fmt.Fprintf(b, "\t%s", column.ColumnName)
		if idx != len(columns)-1 {
			fmt.Fprint(b, ",")
		}
		fmt.Fprintln(b)
	}
	fmt.Fprintf(b, "FROM %s\n", tableName)
	fmt.Fprintf(b, "WHERE %s = ?\n", firstField.ColumnName)
	fmt.Fprint(b, "LIMIT 1;\n")

	return b
}

// SELECT for transactions require some slight changes.
func selectOneSqlTx(tableName string, columns []Column) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	var firstField Column
	if len(columns) > 0 {
		firstField = columns[0]
	}

	fmt.Fprint(b, "SELECT\n")
	for idx, column := range columns {
		fmt.Fprintf(b, "\t%s", column.ColumnName)
		if idx != len(columns)-1 {
			fmt.Fprint(b, ",") // trailing comma except on last line.
		}
		fmt.Fprintln(b)
	}
	fmt.Fprintf(b, "FROM %s\n", tableName)
	fmt.Fprintf(b, "WHERE %s = ?\n", firstField.ColumnName)
	fmt.Fprint(b, "LIMIT 1\n")
	fmt.Fprint(b, "FOR UPDATE;\n")

	return b
}
