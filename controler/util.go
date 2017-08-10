package controler

import (
	"strings"
)

func plural(s string) string {
	if strings.HasSuffix(s, "s") {
		s += "es"
	} else {
		s += "s"
	}

	return s
}
