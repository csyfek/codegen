package mssql

import (
	"bytes"
	"fmt"
)

func  (this *Generator)Delete(typeName, table string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := deleteSql(table, columns)

	funcName := fmt.Sprintf("Delete%s", typeName)
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
	fmt.Fprint(b, b_sql.String())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprintf(b, "\t\t%s, err = db.Prepare(q)", psName)
	fmt.Fprint(b, `
		if err != nil {
		return errors.Stack(err)
		}
	}

	args := []interface{}{id}
	`) // end of prepared statement clause
	fmt.Fprintf(b, "\t_, err = %s.Exec(args...)", psName)
	fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`)

	return b.String()
}

func  (this *Generator)DeleteTx(typeName string, table string, columns []Column) string {

	b := bytes.NewBuffer(nil)
	b_sql := deleteSql(table, columns)

	funcName := fmt.Sprintf("Delete%sTx", typeName)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, id string) error {\n", funcName)
	fmt.Fprintln(b, "q := `")
	fmt.Fprintln(b, b_sql.String())
	fmt.Fprintln(b, "`")

	fmt.Fprint(b, `

	args := []interface{}{id}

	_, err := tx.Exec(q, args...)
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
`)

	return b.String()
}

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func deleteSql(table string, columns []Column) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	fmt.Fprintf(b, "DELETE FROM %s\n", table)
	if len(columns) > 0 {
		column := columns[0]
		fmt.Fprintf(b, "\tWHERE %s.%s = ?;\n", table, column.ColumnName)
	} else {
		fmt.Fprint(b, "\t-- Insert your filter criteria here.\n")
	}

	return b
}
