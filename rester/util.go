package rester

import (
	"github.com/jackmanlabs/codegen"
	"github.com/serenize/snaker"
)

func resource(s string) string {
	//s = strings.ToLower(s)
	s = snaker.CamelToSnake(s)
	s = codegen.Plural(s)

	return s
}
