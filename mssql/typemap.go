package mssql

// TODO: Verify date-time compatibility with time.Time.
var sqlToGo map[string]string = map[string]string{
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

// You'll note here that MSSQL isn't very symmetrical.
var goToSql map[string]string = map[string]string{
	"[]byte":    "varbinary", //
	"bool":      "bit",       //
	"float32":   "real",      // MS-SQL 'real' uses 4 bytes (32 bits).
	"float64":   "float",     // SQL float precision can be variable, but using the max (64-bit) should be safe.
	"int":       "int",       // MS-SQL defines an 'int' to be 32 bits. Go defines it to be 32 or 64 bits. For the sake of convenience, we're simply using 'int'.
	"int16":     "smallint",  // MS-SQL 'smallint' uses 2 bytes (16 bits).
	"int64":     "bigint",    // MS-SQL 'bigint' uses 8 bytes (64 bits).
	"string":    "nvarchar",  //
	"time.Time": "datetime2", //
	"uint":      "tinyint",   // MS-SQL 'tinyint' uses 1 byte (8 bits) and is unsigned.
	//"[]byte":    "binary",     //
	//"float32":   "smallmoney", // MS-SQL does not store currency data.
	//"float64":   "money",      // MS-SQL does not store currency data.
	//"string":    "varchar",    //
	//"time.Time": "date",       //
	//"time.Time": "datetime",   //
	//"time.Time": "time",       //
}

func sqlType(goType string) (string, bool) {
	sqlType, ok := goToSql[goType]
	return sqlType, ok
}
