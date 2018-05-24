package mssql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func (this *generator) SelectOne(pkgName string, def *codegen.Parent) string {

	b := bytes.NewBuffer(nil)
	b_sql := selectOneSql(def)

	funcName := fmt.Sprintf("Get%s", def.Name)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(id string) (*%s.%s, error) {\n", funcName, pkgName, def.Name)
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

	fmt.Fprintf(b, "\tvar x *%s.%s\n", pkgName, def.Name)
	fmt.Fprint(b, "\tif rows.Next() {\n")
	fmt.Fprintf(b, "\t\tx = new(%s.%s)\n", pkgName, def.Name)

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, member := range def.Members {
		fmt.Fprintf(b, "\t\t\t&x.%s,\n", member.GoName)
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

func (this *generator) SelectOneTx(pkgName string, def *codegen.Parent) string {

	b := bytes.NewBuffer(nil)
	b_sql := selectOneSqlTx(def)

	funcName := fmt.Sprintf("Get%sTx", def.Name)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, id string) (*%s.%s, error) {\n", funcName, pkgName, def.Name)
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

	fmt.Fprintf(b, "\tvar x *%s.%s\n", pkgName, def.Name)
	fmt.Fprint(b, "\tif rows.Next() {\n")
	fmt.Fprintf(b, "\t\tx = new(%s.%s)\n", pkgName, def.Name)

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, member := range def.Members {
		fmt.Fprintf(b, "\t\t\t&x.%s,\n", member.GoName)
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

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func selectOneSql(def *codegen.Parent) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	var firstField codegen.Child
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprint(b, "SELECT\n")
	for idx, member := range def.Members {
		if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s.%s\n", def.Table, member.SqlName)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s.%s,\n", def.Table, member.SqlName)
		}
	}
	fmt.Fprintf(b, "FROM %s\n", def.Table)
	fmt.Fprintf(b, "WHERE %s.%s = ?\n", def.Table, firstField.SqlName)
	fmt.Fprint(b, "LIMIT 1;\n")

	return b
}

// SELECT for transactions require some slight changes.
func selectOneSqlTx(def *codegen.Parent) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	var firstField codegen.Child
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprint(b, "SELECT\n")
	for idx, member := range def.Members {
		if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s.%s\n", def.Table, member.SqlName)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s.%s,\n", def.Table, member.SqlName)
		}
	}
	fmt.Fprintf(b, "FROM %s\n", def.Table)
	fmt.Fprintf(b, "WHERE %s.%s = ?\n", def.Table, firstField.SqlName)
	fmt.Fprint(b, "LIMIT 1\n")
	fmt.Fprint(b, "FOR UPDATE;\n")

	return b
}
