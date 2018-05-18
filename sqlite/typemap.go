package sqlite

// TODO: These are not proper mappings, only borrowed from MSSQL as a template.
var sqlToGo map[string]string = map[string]string{
	"BIGINT":     "int64",     // MS-SQL 'bigint' uses 8 bytes (64 bits).
	"BINARY":     "[]byte",    //
	"BIT":        "bool",      //
	"DATE":       "time.Time", //
	"DATETIME":   "time.Time", //
	"DATETIME2":  "time.Time", //
	"FLOAT":      "float64",   // SQL float precision can be variable, but using the max (64-bit) should be safe.
	"INT":        "int",       // MS-SQL defines an 'int' to be 32 bits. Go defines it to be 32 or 64 bits. For the sake of convenience, we're simply using 'int'.
	"MONEY":      "float64",   // MS-SQL does not store currency data.
	"NVARCHAR":   "string",    //
	"REAL":       "float32",   // MS-SQL 'real' uses 4 bytes (32 bits).
	"SMALLINT":   "int16",     // MS-SQL 'smallint' uses 2 bytes (16 bits).
	"SMALLMONEY": "float32",   // MS-SQL does not store currency data.
	"TIME":       "time.Time", //
	"TINYINT":    "uint",      // MS-SQL 'tinyint' uses 1 byte (8 bits) and is unsigned.
	"VARBINARY":  "[]byte",    //
	"VARCHAR":    "string",    //
}

var goToSql map[string]string = map[string]string{
	"[]byte":        "BLOB",
	"bool":          "INTEGER",
	"float32":       "REAL",
	"float64":       "REAL",
	"int":           "INTEGER",
	"int32":         "INTEGER",
	"int64":         "INTEGER",
	"string":        "TEXT",
	"time.Duration": "INTEGER",
	"time.Time":     "DATETIME",
	"uint32":        "INTEGER",
	"uint64":        "INTEGER",
}

func sqlType(goType string) (string, bool) {
	sqlType, ok := goToSql[goType]
	return sqlType, ok
}
