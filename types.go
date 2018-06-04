package codegen

import "github.com/jackmanlabs/tipo"

type Parent struct {
	tipo.Model
	Name     string
	Table    string
	Children []Child
}

type Child struct {
	tipo.Model
	GoName        string // Go-friendly (CamelCase) name.
	SqlName       string // Name to be used for SQL schemas and operations.
	SqlType       string // The SQL type, dependant on the SQL driver.
	SqlConstraint string // Things like 'PRIMARY KEY' and 'NOT NULL'
	JsonName      string // Name for use in JSON-REST APIs.
	GoType        string // All types are normalized to best available Go types.
	Length        int    // Length is preserved for DB-specific configurations.

	// After much debate and trial-and-error, I've settled on keeping the GoName
	// and the SqlName distinct and accessible as strings. Therefore, the caller
	// should be responsible for knowing if and how these values are set and
	// used.
	//
	// The hope is that the GoType and Child objects can be used for translations
	// both to and from SQL schemas/databases.
}
