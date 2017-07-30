package mysql

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/types"
	"strings"
)

// I have to leave out backticks from the SQL because of embedding issues.
// Please refrain from using reserved SQL keywords as struct and member names.
func (this *generator) Schema(pkg *types.Package) string {

	// We need to take enum types and extract their underlying types for the type caster.
	typeMap := make(map[string]*types.Type)
	for _, def := range pkg.Types {
		typeMap[def.Name] = def
	}

	for _, def := range pkg.Types {
		for memberId, member := range def.Members {
			if _, sqlTypeOk := sqlType(member.Type); !sqlTypeOk {
				if t, underlyingTypeOk := typeMap[member.Type]; underlyingTypeOk {
					def.Members[memberId].Type = t.UnderlyingType
				}
			}
		}
	}

	b := bytes.NewBuffer(nil)

	for _, def := range pkg.Types {

		if def.UnderlyingType != "struct" {
			continue
		}

		b.WriteString(this.typeSchema(def))
		fmt.Fprint(b, "\n\n-- -----------------------------------------------------------------------------\n\n")
	}

	return b.String()
}

func (this *generator) typeSchema(def *types.Type) string {
	b := bytes.NewBuffer(nil)

	var firstField types.Member
	if len(def.Members) > 0 {
		firstField = def.Members[0]
	}

	fmt.Fprintf(b, "DROP TABLE IF EXISTS %s;\n\n", def.Table)

	fmt.Fprintf(b, "CREATE TABLE %s (\n", def.Table)

	var columnQty int = 0
	for idx, member := range def.Members {

		var typeSql string
		if member.Type == "string" && (idx == 0 || strings.HasSuffix(member.SqlName, "_id")) {
			// Assume UUID.
			typeSql = "CHAR(36)"
		} else {
			var ok bool
			typeSql, ok = sqlType(member.Type)
			if !ok {
				continue
			}
		}

		columnQty++

		if idx == 0 {
			fmt.Fprintf(b, "\t%s %s PRIMARY KEY,\n", member.SqlName, typeSql)
		} else if idx == len(def.Members)-1 {
			fmt.Fprintf(b, "\t%s %s NOT NULL\n", member.SqlName, typeSql)
		} else {
			fmt.Fprintf(b, "\t%s %s NOT NULL,\n", member.SqlName, typeSql)
		}
	}

	if columnQty == 0{
		return ""
	}

	fmt.Fprintf(b, "\t-- FOREIGN KEY (%s) REFERENCES parent_table (id) ON DELETE CASCADE\n", firstField.SqlName)
	fmt.Fprint(b, ")\n")
	fmt.Fprint(b, "\t-- ENGINE = TokuDB\n")
	fmt.Fprint(b, "\tENGINE = InnoDB\n")
	fmt.Fprint(b, "\tROW_FORMAT = COMPRESSED -- Requires Barracuda file format.\n")
	fmt.Fprint(b, "\tDEFAULT CHARSET = utf8;\n")

	return b.String()

}
