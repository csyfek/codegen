package pg

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/extractor"
	"github.com/serenize/snaker"
)

func SelectPlural(pkgName string, def *extractor.StructDefinition) string {

	members := getGoSqlData(def.Members)

	b := bytes.NewBuffer(nil)
	b_sql := selectPluralSql(def, members)

	funcName := fmt.Sprintf("Get%ss", def.Name)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s() ([]%s.%s, error) {\n", funcName, pkgName, def.Name)
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
	fmt.Fprint(b, "\targs := []interface{}{}\n\n")
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
	for _, member := range members {
		if !member.SqlCompatible {
			fmt.Fprintf(b, "\t\tvar x_%s []byte\n", member.Name)
		}
	}

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, member := range members {
		if member.SqlCompatible {
			fmt.Fprintf(b, "\t\t\t&x.%s,\n", member.Name)
		} else {
			fmt.Fprintf(b, "\t\t\t&x_%s,\n", member.Name)
		}
	}

	fmt.Fprint(b, "\t\t}\n") // end of targets declaration.
	fmt.Fprint(b, `
		err = rows.Scan(targets...)
		if err != nil {
			return z, errors.Stack(err)
		}

`)

	for _, member := range members {
		if !member.SqlCompatible {
			fmt.Fprintf(b, "\t\terr = json.Unmarshal(x_%s, &x.%s)", member.Name, member.Name)
			fmt.Fprint(b, `
		if err != nil {
			return z, errors.Stack(err)
		}

`)
		}
	}

	fmt.Fprint(b, "\t\tz = append(z, x)\n")
	fmt.Fprint(b, "\t}\n\n") // end of scan clause.
	fmt.Fprint(b, "\t// empty slice is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn z, nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

func SelectPluralTx(pkgName string, def *extractor.StructDefinition) string {

	members := getGoSqlData(def.Members)

	b := bytes.NewBuffer(nil)
	b_sql := selectPluralSql(def, members)

	funcName := fmt.Sprintf("Get%ssTx", def.Name)

	fmt.Fprintf(b, "func %s(tx *sql.Tx) ([]%s.%s, error) {\n", funcName, pkgName, def.Name)

	fmt.Fprint(b, "\t\tq := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n\n")

	fmt.Fprint(b, "\targs := []interface{}{}\n\n")
	fmt.Fprint(b, "\trows, err := tx.Query(q, args...)")
	fmt.Fprint(b, `
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer rows.Close()

`)

	fmt.Fprintf(b, "\tvar z []%s.%s = make([]%s.%s, 0)\n", pkgName, def.Name, pkgName, def.Name)
	fmt.Fprint(b, "\tfor rows.Next() {\n")
	fmt.Fprintf(b, "\t\tvar x %s.%s\n", pkgName, def.Name)
	for _, member := range members {
		if !member.SqlCompatible {
			fmt.Fprintf(b, "\t\tvar x_%s []byte\n", member.Name)
		}
	}

	fmt.Fprint(b, "\t\ttargets := []interface{}{\n")
	for _, member := range members {
		if member.SqlCompatible {
			fmt.Fprintf(b, "\t\t\t&x.%s,\n", member.Name)
		} else {
			fmt.Fprintf(b, "\t\t\t&x_%s,\n", member.Name)
		}
	}

	fmt.Fprint(b, "\t\t}\n") // end of targets declaration.
	fmt.Fprint(b, `
		err = rows.Scan(targets...)
		if err != nil {
			return z, errors.Stack(err)
		}

`)

	for _, member := range members {
		if !member.SqlCompatible {
			fmt.Fprintf(b, "\t\terr = json.Unmarshal(x_%s, &x.%s)", member.Name, member.Name)
			fmt.Fprint(b, `
		if err != nil {
			return z, errors.Stack(err)
		}

`)
		}
	}

	fmt.Fprint(b, "\t\tz = append(z, x)\n")
	fmt.Fprint(b, "\t}\n\n") // end of scan clause.
	fmt.Fprint(b, "\t// empty slice is returned if no data was present.\n")
	fmt.Fprint(b, "\treturn z, nil\n")

	fmt.Fprint(b, "}\n") // end of function

	return b.String()
}

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func selectPluralSql(def *extractor.StructDefinition, members []GoSqlDatum) *bytes.Buffer {

	b := bytes.NewBuffer(nil)
	tableName := snaker.CamelToSnake(def.Name)

	var firstField GoSqlDatum
	if len(members) > 0 {
		firstField = members[0]
	}

	fmt.Fprint(b, "SELECT\n")
	for idx, member := range members {
		if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s.%s\n", tableName, member.SqlName)
		} else {
			// Note the trailing comma.
			fmt.Fprintf(b, "\t%s.%s,\n", tableName, member.SqlName)
		}
	}
	fmt.Fprintf(b, "FROM %s;\n", tableName)
	fmt.Fprint(b, "-- Update your filter criteria here:\n")
	fmt.Fprintf(b, "-- WHERE %s.%s = $1;\n", tableName, firstField.SqlName)

	return b
}
