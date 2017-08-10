package mssql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/common"
)

func (this *generator) Delete(def *common.Type) string {

	b := bytes.NewBuffer(nil)
	b_sql := deleteSql(def)

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

func (this *generator) DeleteTx(def *common.Type) string {

	b := bytes.NewBuffer(nil)
	b_sql := deleteSql(def)

	funcName := fmt.Sprintf("Delete%sTx", def.Name)

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
func deleteSql(def *common.Type) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	fmt.Fprintf(b, "DELETE FROM %s\n", def.Table)
	if len(def.Members) > 0 {
		column := def.Members[0]
		fmt.Fprintf(b, "\tWHERE %s.%s = ?;\n", def.Table, column.SqlName)
	} else {
		fmt.Fprint(b, "\t-- Insert your filter criteria here.\n")
	}

	return b
}
