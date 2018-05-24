package pg

import "github.com/jackmanlabs/codegen"

type Column struct {
	TableCatalog           string
	TableSchema            string
	TableName              string
	ColumnName             string
	OrdinalPosition        int
	ColumnDefault          *string
	IsNullable             string
	DataType               string
	CharacterMaximumLength *int
	CharacterOctetLength   *int
	NumericPrecision       *int
	NumericPrecisionRadix  *int
	NumericScale           *int
	DatetimePrecision      *int
	CharacterSetCatalog    *string
	CharacterSetSchema     *string
	CharacterSetName       *string
	CollationCatalog       *string
	CollationSchema        *string
	CollationName          *string
	DomainCatalog          *string
	DomainSchema           *string
	DomainName             *string
}

func (this *Column) Member() codegen.Child {

	var l int

	if this.CharacterMaximumLength != nil {
		l = *this.CharacterMaximumLength
	} else if this.NumericPrecision != nil {
		l = *this.NumericPrecision
	}

	return codegen.Child{
		// We expect the Go and SQL names to be the same for MSSQL.
		GoName:  this.ColumnName,
		SqlName: this.ColumnName,
		GoType:  this.goType(),
		Length:  l,
	}
}

func (this *Column) goType() string {

	for s, g := range sqlToGo {
		if s == this.DataType {
			return g
		}
	}

	return "mysql_" + this.DataType
}
