package mysql

var goToSql map[string]string = map[string]string{
	"[]byte":        "BLOB",
	"bool":          "TINYINT",
	"float32":       "DOUBLE",
	"float64":       "DOUBLE",
	"int":           "INT",
	"int32":         "INT",
	"int64":         "BIGINT",
	"string":        "VARCHAR(255)",
	"time.Duration": "BIGINT",
	"time.Time":     "DATETIME",
	"uint32":        "INT UNSIGNED",
	"uint64":        "BIGINT UNSIGNED",
}

var sqlToGo map[string]string = map[string]string{
	"BIGINT UNSIGNED": "uint64",
	"BIGINT":          "int64",
	"BLOB":            "[]byte",
	"DATETIME":        "time.Time",
	"DOUBLE":          "float64",
	"INT UNSIGNED":    "uint32",
	"INT":             "int",
	"TINYINT":         "bool",
	"VARCHAR(255)":    "string",
}

func sqlType(goType string) (string, bool) {
	sqlType, ok := goToSql[goType]
	return sqlType, ok
}
