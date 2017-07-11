package sqlite

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/types"
	"github.com/serenize/snaker"
)

func (this *generator) UpdateOne(pkgName string, def *types.Type) string {

	b := bytes.NewBuffer(nil)
	b_sql := updateSql(def)

	funcName := fmt.Sprintf("Update%s", def.Name)
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

	for _, member := range def.Members {
		if _, ok := sqlType(member.Type); !ok {
			fmt.Fprintf(b, "\tvar x_%s []byte\n", member.GoName)
		}
	}
	fmt.Fprint(b, "\n")

	for _, member := range def.Members {
		if _, ok := sqlType(member.Type); !ok {
			fmt.Fprintf(b, "\tx_%s, err = json.Marshal(x.%s)", member.GoName, member.GoName)
			fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

`)
		}
	}

	fmt.Fprint(b, "\targs := []interface{}{\n")
	for _, member := range def.Members {
		if _, ok := sqlType(member.Type); ok {
			fmt.Fprintf(b, "\t\t&x.%s,\n", member.GoName)
		} else {
			fmt.Fprintf(b, "\t\t&x_%s,\n", member.GoName)
		}
	}
	if len(def.Members) > 0 {
		fmt.Fprintf(b, "\t\t&x.%s,\n", def.Members[0].GoName)
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

func (this *generator) UpdateOneTx(pkgName string, def *types.Type) string {

	b := bytes.NewBuffer(nil)
	b_sql := updateSql(def)

	funcName := fmt.Sprintf("Update%sTx", def.Name)

	fmt.Fprintf(b, "func %s(tx *sql.Tx, x *%s.%s) error {\n", funcName, pkgName, def.Name)
	fmt.Fprint(b, "var err error\n")
	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	for _, member := range def.Members {
		if _, ok := sqlType(member.Type); !ok {
			fmt.Fprintf(b, "\tvar x_%s []byte\n", member.GoName)
		}
	}
	fmt.Fprint(b, "\n")

	for _, member := range def.Members {
		if _, ok := sqlType(member.Type); !ok {
			fmt.Fprintf(b, "\tx_%s, err = json.Marshal(x.%s)", member.GoName, member.GoName)
			fmt.Fprint(b, `
	if err != nil {
		return errors.Stack(err)
	}

`)
		}
	}

	fmt.Fprint(b, "\targs := []interface{}{\n")
	for _, member := range def.Members {
		if _, ok := sqlType(member.Type); ok {
			fmt.Fprintf(b, "\t\t&x.%s,\n", member.GoName)
		} else {
			fmt.Fprintf(b, "\t\t&x_%s,\n", member.GoName)
		}
	}
	if len(def.Members) > 0 {
		fmt.Fprintf(b, "\t\t&x.%s,\n", def.Members[0].GoName)
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

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func updateSql(def *types.Type) *bytes.Buffer {

	b := bytes.NewBuffer(nil)
	tableName := snaker.CamelToSnake(def.Name)

	var firstField types.Member
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprintf(b, "UPDATE %s\n", tableName)
	fmt.Fprint(b, "SET\n")
	for idx, member := range def.Members {
		if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s = $%d\n", member.SqlName, idx+1)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s = $%d,\n", member.SqlName, idx+1)
		}
	}
	fmt.Fprintf(b, "WHERE %s.%s = $%d;\n", tableName, firstField.SqlName, len(def.Members)+1)

	return b
}
