package pkger

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/common"
	"github.com/segmentio/go-camelcase"
)

// https://google.github.io/styleguide/jsoncstyleguide.xml#Property_Name_Format

func GenerateModel(def *common.Type) (string, []string) {
	var (
		b       *bytes.Buffer = bytes.NewBuffer(nil)
		imports []string      = make([]string, 0)
	)

	fmt.Fprintf(b, "type %s struct{\n", def.Name)

	for _, member := range def.Members {

		// We assume the Go Name is PascalCase.
		//jsonName := snaker.SnakeToCamelLower(member.GoName)
		jsonName := camelcase.Camelcase(member.GoName)

		fmt.Fprintf(b, "\t%s\t%s\t`json:\"%s\"`\n", member.GoName, member.GoType, jsonName)

		if member.GoType == "time.Time" && !sContains(imports, "time") {
			imports = append(imports, "time")
		}
	}

	fmt.Fprint(b, "}\n")

	return b.String(), imports
}
