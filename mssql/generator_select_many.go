package mssql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen"
)

func (this *generator) SelectMany(pkgName string, def *codegen.Parent) string {

	b := bytes.NewBuffer(nil)
	b_sql := selectManySql(def)

	funcName := fmt.Sprintf("Get%ss", def.Name)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(filter filters.%s) ([]%s.%s, error) {\n", funcName, def.Name, pkgName, def.Name)
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
	fmt.Fprint(b, `
	args := []interface{}{
		// If you add a filter criteria, insert it here:
		// filter,
	}

	`)
	fmt.Fprintf(b, "\trows, err := %s.Query(args...)", psName)
	fmt.Fprint(b, `
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

`)

	fmt.Fprintf(b, "\tvar z []%s.%s = make([]%s.%s, 0)\n", pkgName, def.Name, pkgName, def.Name)
	fmt.Fprint(b, "\tfor rows.Next() {\n")
	fmt.Fprintf(b, "\t\tvar x %s.%s\n", pkgName, def.Name)

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, column := range def.Members {
		fmt.Fprintf(b, "\t\t\t&x.%s,\n", column.GoName)
	}

	fmt.Fprint(b, "\t\t}\n") // end of targets declaration.
	fmt.Fprint(b, `
		err = rows.Scan(targets...)
		if err != nil {
			return z, errors.Stack(err)
		}

`)

	fmt.Fprint(b, "\t\tz = append(z, x)\n")
	fmt.Fprint(b, "\t}\n\n") // end of scan clause.
	fmt.Fprint(b, "\t// empty slice is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn z, nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

func (this *generator) SelectManyTx(pkgName string, def *codegen.Parent) string {

	b := bytes.NewBuffer(nil)
	b_sql := selectManySqlTx(def)

	funcName := fmt.Sprintf("Get%ssTx", def.Name)

	fmt.Fprintf(b, "func %s(tx *sql.Tx ,filter filters.%s) ([]%s.%s, error) {\n", funcName, def.Name, pkgName, def.Name)

	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`") // backtick needs to be in double quotes.

	fmt.Fprint(b, `

	args := []interface{}{
		// If you add a filter criteria, insert it here:
		// filter,
	}

	rows, err := tx.Query(q, args...)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

`)

	fmt.Fprintf(b, "\tvar z []%s.%s = make([]%s.%s, 0)\n", pkgName, def.Name, pkgName, def.Name)
	fmt.Fprint(b, "\tfor rows.Next() {\n")
	fmt.Fprintf(b, "\t\tvar x %s.%s\n", pkgName, def.Name)

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, column := range def.Members {
		fmt.Fprintf(b, "\t\t\t&x.%s,\n", column.GoName)
	}

	fmt.Fprint(b, `
		} // end of targets declaration.

		err = rows.Scan(targets...)
		if err != nil {
			return z, errors.Stack(err)
		}

`)

	fmt.Fprint(b, `
		z = append(z, x)
	} // end of scan clause.

	// empty slice is returned if no data was present.
	return z, nil
}
`)
	return b.String()
}

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and column names.
func selectManySql(def *codegen.Parent) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	fmt.Fprint(b, "SELECT\n")
	for idx, column := range def.Members {
		fmt.Fprintf(b, "\t%s", column.SqlName)
		if idx != len(def.Members)-1 {
			fmt.Fprint(b, ",") // trailing comma except on last line.
		}
		fmt.Fprintln(b)
	}
	fmt.Fprintf(b, "FROM %s;\n", def.Table)
	fmt.Fprint(b, "-- Update your filter criteria here:\n")
	fmt.Fprint(b, "-- WHERE filter = ?;\n")

	return b
}

// SELECT for transactions require some slight changes.
func selectManySqlTx(def *codegen.Parent) *bytes.Buffer {

	b := bytes.NewBuffer(nil)

	fmt.Fprint(b, "SELECT\n")
	for idx, column := range def.Members {
		fmt.Fprintf(b, "\t%s", column.SqlName)
		if idx != len(def.Members)-1 {
			fmt.Fprint(b, ",") // trailing comma except on last line.
		}
		fmt.Fprintln(b)
	}
	fmt.Fprintf(b, "FROM %s\n", def.Table)
	fmt.Fprint(b, "-- Update your filter criteria here:\n")
	fmt.Fprint(b, "-- WHERE filter = ?\n")
	fmt.Fprint(b, "LIMIT 1\n")
	fmt.Fprint(b, "FOR UPDATE;\n")

	return b
}
