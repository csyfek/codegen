package pg

var sqlToGo map[string]string = map[string]string{
	"[]byte":        "BLOB",
	"bool":          "TINYINT",
	"float32":       "DOUBLE PRECISION",
	"float64":       "DOUBLE PRECISION",
	"int":           "INT",
	"int32":         "INT",
	"int64":         "BIGINT",
	"string":        "VARCHAR(255)",
	"time.Duration": "BIGINT",
	"time.Time":     "TIMESTAMP",
	"uint32":        "INT UNSIGNED",
	"uint64":        "BIGINT UNSIGNED",
}

func sqlType(goType string) (string, bool) {
	sqlType, ok := goToSql[goType]
	return sqlType, ok
}

var goToSql map[string]string = map[string]string{
	"int":           "INT",
	"string":        "VARCHAR(255)",
	"float32":       "DOUBLE PRECISION",
	"float64":       "DOUBLE PRECISION",
	"time.Time":     "TIMESTAMP",
	"time.Duration": "BIGINT",
	"int64":         "BIGINT",
	"int32":         "INT",
	"uint32":        "INT UNSIGNED",
	"uint64":        "BIGINT UNSIGNED",
	"bool":          "TINYINT",
	"[]byte":        "BLOB",
}
