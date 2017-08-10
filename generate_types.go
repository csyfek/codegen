package main

import (
	"fmt"
	"github.com/jackmanlabs/codegen/common"
	"github.com/jackmanlabs/codegen/pkger"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"io"
	"os"
)

func generateTypes(outputRoot string, pkg *common.Package) error {

	var (
		f    io.WriteCloser
		path string
	)

	path = outputRoot + "/types"
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

		fmt.Fprintf(f, "package %s\n\n", pkg.Name)

		strct, importPaths := pkger.GenerateStruct(def)

		fmt.Fprint(f, "import(\n")
		for _, importPath := range importPaths {
			fmt.Fprintf(f, "\t\"%s\"\n", importPath)
		}
		fmt.Fprint(f, ")\n\n")

		fmt.Fprint(f, strct)

		err = f.Close()
		if err != nil {
			return errors.Stack(err)
		}
	}

	return nil
}
