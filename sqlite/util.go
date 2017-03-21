package sqlite

import (
	"github.com/jackmanlabs/codegen/extractor"
	"github.com/serenize/snaker"
)

func getSqlType(goType string) (sqlType string, sqlCompatible bool) {
	switch goType {
	case "[]byte":
		return "BLOB", true
	case "bool":
		return "INTEGER", true
	case "float32":
		return "REAL", true
	case "float64":
		return "REAL", true
	case "int":
		return "INTEGER", true
	case "int32":
		return "INTEGER", true
	case "int64":
		return "INTEGER", true
	case "string":
		return "TEXT", true
	case "time.Duration":
		return "INTEGER", true
	case "time.Time":
		return "TEXT", true
	case "uint32":
		return "INTEGER", true
	case "uint64":
		return "INTEGER", true // 64-bit
	}

	return "TEXT", false
}

type GoSqlDatum struct {
	extractor.StructMemberDefinition
	SqlCompatible bool
	SqlType       string
	SqlName       string
}

func getGoSqlData(structMembers []extractor.StructMemberDefinition) []GoSqlDatum {
	members := make([]GoSqlDatum, 0)
	for _, member_ := range structMembers {
		sqlType, compatible := getSqlType(member_.Type)
		member := GoSqlDatum{
			StructMemberDefinition: member_,
			SqlType:                sqlType,
			SqlCompatible:          compatible,
			SqlName:                snaker.CamelToSnake(member_.Name),
		}
		members = append(members, member)
	}

	return members
}
