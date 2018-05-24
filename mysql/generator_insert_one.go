package mysql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func (this *generator) InsertOne(pkgName string, def *codegen.Parent) string {

	b := bytes.NewBuffer(nil)
	b_sql := insertSql(def)

	funcName := fmt.Sprintf("Insert%s", def.Name)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(x *%s.%s) error {\n", funcName, pkgName, def.Name)
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

	fmt.Fprint(b, "\targs := []interface{}{\n")
	for _, member := range def.Members {
		fmt.Fprintf(b, "\t\t&x.%s,\n", member.GoName)
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

func (this *generator) InsertOneTx(pkgName string, def *codegen.Parent) string {

	b := bytes.NewBuffer(nil)
	b_sql := insertSql(def)

	funcName := fmt.Sprintf("Insert%sTx", def.Name)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, x *%s.%s) error {\n", funcName, pkgName, def.Name)
	fmt.Fprint(b, "var err error\n")

	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprint(b, "\targs := []interface{}{\n")
	for _, member := range def.Members {
		fmt.Fprintf(b, "\t\t&x.%s,\n", member.GoName)
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
func insertSql(def *codegen.Parent) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	fmt.Fprintf(b, "INSERT INTO %s (\n", def.Table)
	for idx, member := range def.Members {
		if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s\n", member.SqlName)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s,\n", member.SqlName)
		}
	}
	fmt.Fprint(b, ") VALUES (")
	for idx := range def.Members {
		if idx == len(def.Members)-1 {
			fmt.Fprint(b, "?);\n")
		} else {
			fmt.Fprint(b, "?, ")
		}
	}

	return b
}
