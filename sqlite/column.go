package sqlite

import "github.com/jackmanlabs/codegen/types"

// TODO: This column is borrowed from MSSQL as a template.
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

func (this *Column) Member() types.Member {

	var l int

	if this.CharacterMaximumLength != nil {
		l = *this.CharacterMaximumLength
	} else if this.NumericPrecision != nil {
		l = *this.NumericPrecision
	}

	return types.Member{
		GoName: this.ColumnName,
		Type:   this.goType(),
		Length: l,
	}
}

func (this *Column) goType() string {

	for s, g := range sqlTogo {
		if s == this.DataType {
			return g
		}
	}

	return "sqlite_" + this.DataType

}
