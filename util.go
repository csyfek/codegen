package codegen

import (
	"strings"
)

func Plural(s string) string {
	switch {
	case strings.HasSuffix(s, "y"):
		strings.TrimSuffix(s, "y")
		s += "ies"
	case strings.HasSuffix(s, "s"):
		s += "es"
	default:
		s += "s"
	}

	return s
}
