package main

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/common"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"io"
	"os"
)

// This is a bit of a hack at the moment. I've never seen explicit filter types
// declared or used, and I want to see if it's a dead end.

func generateFilters(outputRoot string, pkg *common.Package) error {

	var (
		f    io.WriteCloser
		path string
	)

	path = outputRoot + "/filters"
	d, err := os.Open(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir|os.ModePerm)
		if err != nil {
			return errors.Stack(err)
		}
	} else if err != nil {
		return errors.Stack(err)
	} else {
		err = d.Close()
		if err != nil {
			return errors.Stack(err)
		}
	}

	for _, def := range pkg.Types {

		if def.UnderlyingType != "struct" {
			continue
		}

		filename := snaker.CamelToSnake(def.Name) + ".go"

		f, err = os.Create(path + "/" + filename)
		if err != nil {
			return errors.Stack(err)
		}

		fmt.Fprintf(f, "package %s\n\n", "filters")

		filter, importPaths := generateFilter(def)

		fmt.Fprint(f, "import(\n")
		for _, importPath := range importPaths {
			fmt.Fprintf(f, "\t\"%s\"\n", importPath)
		}
		fmt.Fprint(f, ")\n\n")

		fmt.Fprint(f, filter)

		err = f.Close()
		if err != nil {
			return errors.Stack(err)
		}
	}

	return nil
}

func generateFilter(def *common.Type) (string, []string) {
	var (
		b       *bytes.Buffer = bytes.NewBuffer(nil)
		imports []string      = make([]string, 0)
	)

	fmt.Fprintf(b, "type %s struct{\n", def.Name)

	for _, member := range def.Members {

		goType := "*" + member.GoType

		if member.IsNumeric() {
			fmt.Fprintf(b, "\t%s\t%s\n", member.GoName, goType)
			fmt.Fprintf(b, "\t%s_Min\t%s\n", member.GoName, goType)
			fmt.Fprintf(b, "\t%s_Max\t%s\n", member.GoName, goType)
		} else {
			fmt.Fprintf(b, "\t%s\t%s\n", member.GoName, goType)
		}

		if member.GoType == "time.Time" && !sContains(imports, "time") {
			imports = append(imports, "time")
		}

	}

	fmt.Fprint(b, "}\n")

	return b.String(), imports
}

func sContains(set []string, s string) bool {
	for _, s_ := range set {
		if s == s_ {
			return true
		}
	}
	return false
}
