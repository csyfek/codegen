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
		Name:   this.ColumnName,
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

// TODO: These are not proper mappings, only borrowed from MSSQL as a template.
var sqlTogo map[string]string = map[string]string{
	"bigint":     "int64",     // MS-SQL 'bigint' uses 8 bytes (64 bits).
	"binary":     "[]byte",    //
	"bit":        "bool",      //
	"date":       "time.Time", //
	"datetime":   "time.Time", //
	"datetime2":  "time.Time", //
	"float":      "float64",   // SQL float precision can be variable, but using the max (64-bit) should be safe.
	"int":        "int",       // MS-SQL defines an 'int' to be 32 bits. Go defines it to be 32 or 64 bits. For the sake of convenience, we're simply using 'int'.
	"money":      "float64",   // MS-SQL does not store currency data.
	"nvarchar":   "string",    //
	"real":       "float32",   // MS-SQL 'real' uses 4 bytes (32 bits).
	"smallint":   "int16",     // MS-SQL 'smallint' uses 2 bytes (16 bits).
	"smallmoney": "float32",   // MS-SQL does not store currency data.
	"time":       "time.Time", //
	"tinyint":    "uint",      // MS-SQL 'tinyint' uses 1 byte (8 bits) and is unsigned.
	"varbinary":  "[]byte",    //
	"varchar":    "string",    //
}
