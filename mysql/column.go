package mysql

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

func (this *Column) SqlType() string {
	return this.DataType
}

func (this *Column) SqlName() string {
	return this.ColumnName
}

func (this *Column) GoType() string {

	var t string

	switch this.DataType {
	case "bigint":
		// MS-SQL 'bigint' uses 8 bytes (64 bits).
		t = "int64"
	case "binary":
		t = "[]byte"
	case "bit":
		t = "bool"
	case "date":
		t = "time.Time"
	case "datetime":
		t = "time.Time"
	case "datetime2":
		t = "time.Time"
	case "float":
		t = "float64"
		if this.NumericPrecision != nil && *this.NumericPrecision < 24 {
			t = "float32"
		}
	case "int":
		// MS-SQL defines an 'int' to be 32 bits. Go defines it to be 32 or 64 bits.
		// For the sake of convenience, we're simply using 'int'.
		t = "int"
	case "money":
		// MS-SQL does not store currency data.
		t = "float64"
	case "nvarchar":
		t = "string"
	case "real":
		// MS-SQL 'real' uses 4 bytes (32 bits).
		t = "float32"
	case "smallint":
		// MS-SQL 'smallint' uses 2 bytes (16 bits).
		t = "int16"
	case "smallmoney":
		// MS-SQL does not store currency data.
		t = "float32"
	case "time":
		t = "time.Time"
	case "tinyint":
		// MS-SQL 'tinyint' uses 1 byte (8 bits) and is unsigned.
		t = "uint"
	case "varbinary":
		t = "[]byte"
	case "varchar":
		t = "string"
	default:
		t = "mssql_" + this.DataType
	}

	return t
}
