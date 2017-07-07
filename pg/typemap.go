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
