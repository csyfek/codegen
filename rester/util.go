package rester

import (
	"github.com/serenize/snaker"
	"strings"
)

func resource(s string) string {
	//s = strings.ToLower(s)
	s = snaker.CamelToSnake(s)
	s = plural(s)

	return s
}

func plural(s string) string {
	if strings.HasSuffix(s, "s") {
		s += "es"
	} else {
		s += "s"
	}

	return s
}
