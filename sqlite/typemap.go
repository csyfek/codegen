package sqlite

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
	"time.Time":     "TEXT",
	"uint32":        "INTEGER",
	"uint64":        "INTEGER",
}

func sqlType(goType string) (string, bool) {
	sqlType, ok := goToSql[goType]
	return sqlType, ok
}
