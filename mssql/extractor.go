package mssql

import (
	"database/sql"
	"fmt"
	"github.com/jackmanlabs/codegen/types"
	"github.com/jackmanlabs/errors"
	"strings"
	"sync"
)

func (this *Extractor) db() (*sql.DB, error) {
	this.Lock()
	defer this.Unlock()

	if this._db != nil {
		return this._db, nil
	}

	connString := fmt.Sprintf("server=%s;database=%s;user id=%s;password=%s", this.hostname, this.database, this.username, this.password)

	var err error
	this._db, err = sql.Open("mssql", connString)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return this._db, nil
}

type Extractor struct {
	_db *sql.DB
	sync.Mutex
	username string
	password string
	hostname string
	database string
}

func NewExtractor(username, password, hostname, database string) *Extractor {
	this := &Extractor{
		username: username,
		password: password,
		hostname: hostname,
		database: database,
	}
	return this
}

func (this *Extractor) Extract() (*types.Package, error) {

	tables, err := this.tables()
	if err != nil {
		return nil, errors.Stack(err)
	}

	tableColumns := make(map[string][]Column)

	for _, table := range tables {
		if !strings.HasPrefix(table, "tbl") {
			continue
		}

		columns, err := this.columns(table)
		if err != nil {
			return nil, errors.Stack(err)
		}

		tableColumns[table] = columns
	}

	pkg := &types.Package{
		Types:   make([]*types.Type, 0),
		Imports: nil,
		Name:    "",
		Path:    "",
	}

	for table, columns := range tableColumns {

		t := types.NewType()
		t.Name = strings.TrimPrefix(table, "tbl")

		for _, column := range columns {
			t.Members = append(t.Members, column.Member())
		}
	}

	return pkg, nil
}

func (this *Extractor) tables() ([]string, error) {

	db, err := this.db()
	if err != nil {
		return nil, errors.Stack(err)
	}

	q := `SELECT Distinct TABLE_NAME FROM information_schema.TABLES;`

	rows, err := db.Query(q)
	if err != nil {
		return nil, errors.Stack(err)
	}

	tables := make([]string, 0)
	for rows.Next() {
		var s string

		err = rows.Scan(&s)
		if err != nil {
			return nil, errors.Stack(err)
		}

		tables = append(tables, s)
	}

	return tables, nil
}

func (this *Extractor) columns(table string) ([]Column, error) {

	db, err := this.db()
	if err != nil {
		return nil, errors.Stack(err)
	}

	q := `
SELECT
  TABLE_CATALOG,
  TABLE_SCHEMA,
  TABLE_NAME,
  COLUMN_NAME,
  ORDINAL_POSITION,
  COLUMN_DEFAULT,
  IS_NULLABLE,
  DATA_TYPE,
  CHARACTER_MAXIMUM_LENGTH,
  CHARACTER_OCTET_LENGTH,
  NUMERIC_PRECISION,
  NUMERIC_PRECISION_RADIX,
  NUMERIC_SCALE,
  DATETIME_PRECISION,
  CHARACTER_SET_CATALOG,
  CHARACTER_SET_SCHEMA,
  CHARACTER_SET_NAME,
  COLLATION_CATALOG,
  COLLATION_SCHEMA,
  COLLATION_NAME,
  DOMAIN_CATALOG,
  DOMAIN_SCHEMA,
  DOMAIN_NAME
FROM INFORMATION_SCHEMA.COLUMNS
WHERE TABLE_NAME = ?;
`

	rows, err := db.Query(q, table)
	if err != nil {
		return nil, errors.Stack(err)
	}

	columns := make([]Column, 0)
	for rows.Next() {
		var c Column

		targets := []interface{}{
			&c.TableCatalog,
			&c.TableSchema,
			&c.TableName,
			&c.ColumnName,
			&c.OrdinalPosition,
			&c.ColumnDefault,
			&c.IsNullable,
			&c.DataType,
			&c.CharacterMaximumLength,
			&c.CharacterOctetLength,
			&c.NumericPrecision,
			&c.NumericPrecisionRadix,
			&c.NumericScale,
			&c.DatetimePrecision,
			&c.CharacterSetCatalog,
			&c.CharacterSetSchema,
			&c.CharacterSetName,
			&c.CollationCatalog,
			&c.CollationSchema,
			&c.CollationName,
			&c.DomainCatalog,
			&c.DomainSchema,
			&c.DomainName,
		}

		err = rows.Scan(targets...)
		if err != nil {
			return nil, errors.Stack(err)
		}
		columns = append(columns, c)
	}

	return columns, nil
}
