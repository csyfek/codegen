package generate_sql

import (
	"github.com/jackmanlabs/codegen/structfinder"
	"github.com/serenize/snaker"
)

func getSqlType(goType string) (sqlType string, sqlCompatible bool) {
	switch goType {
	case "string":
		return "VARCHAR(256)", true
	case "float32":
		return "DOUBLE", true
	case "float64":
		return "DOUBLE", true
	case "time.Time":
		return "DATETIME", true
	case "int64":
		return "BIGINT", true
	case "int32":
		return "INT", true
	case "uint32":
		return "INT UNSIGNED", true
	case "uint64":
		return "BIGINT UNSIGNED", true // 64-bit
	case "bool":
		return "TINYINT", true
	case "[]byte":
		return "BLOB", true
	}

	return "BLOB", false
}

type GoSqlDatum struct {
	structfinder.StructMemberDefinition
	SqlCompatible bool
	SqlType       string
	SqlName       string
}


func getGoSqlData(structMembers []structfinder.StructMemberDefinition) []GoSqlDatum {
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
