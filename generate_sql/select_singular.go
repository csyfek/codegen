package generate_sql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/structfinder"
)

func SelectSingular(def structfinder.StructDefinition) string {

	b := bytes.NewBuffer(nil)
	b_sql := bytes.NewBuffer(nil)

	funcName := fmt.Sprintf("Get%s", def.Name)
	psName := fmt.Sprintf("ps_%s", funcName)

	fmt.Fprintf(b, "var %s *sql.Stmt\n\n", psName)
	fmt.Fprintf(b, "func %s(id string) (%s.%s, error) {\n", funcName, def.Package, def.Name)
	fmt.Fprintf(b, "\tif %s != nil{\n", psName)
	fmt.Fprint(b, "\t\tsql := `\n")
	fmt.Fprintf(b, "%s", b_sql.Bytes())
	fmt.Fprint(b, "`\n")

	fmt.Fprint(b, "	}\n")
	fmt.Fprint(b, "}\n")

	return b.String()
}

func generateSelectPlural(structfinder.StructDefinition) string {

	return ""
}
