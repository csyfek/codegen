package rester

import (
	"github.com/serenize/snaker"
	"github.com/jackmanlabs/codegen"
)

func resource(s string) string {
	//s = strings.ToLower(s)
	s = snaker.CamelToSnake(s)
	s = codegen.Plural(s)

	return s
}

