package main

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/common"
	"github.com/jackmanlabs/codegen/rester"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"io"
	"os"
	"sort"
)

func generateRest(outputRoot string, importPaths []string, pkg *common.Package) error {

	path := outputRoot

	registers := make([]string, 0)

	for _, def := range pkg.Types {

		b := bytes.NewBuffer(nil)

		fmt.Fprint(b, `
package main

import(
	"strconv"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jackmanlabs/errors"
`)

		for _, importPath := range importPaths {
			if importPath != "" {
				fmt.Fprintf(b, "\t\"%s\"\n", importPath)
			}
		}
		fmt.Fprint(b, ")\n\n")

		var register, handler string

		fmt.Fprintln(b)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b, "// TYPE: "+def.Name)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b)

		register, handler = rester.GetOne(def)
		registers = append(registers, register)
		fmt.Fprintln(b, handler)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		register, handler = rester.GetCollection(def)
		registers = append(registers, register)
		fmt.Fprintln(b, handler)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		register, handler = rester.PostCollection(def)
		registers = append(registers, register)
		fmt.Fprintln(b, handler)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		register, handler = rester.PutOne(def)
		registers = append(registers, register)
		fmt.Fprintln(b, handler)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		register, handler = rester.DeleteOne(def)
		registers = append(registers, register)
		fmt.Fprintln(b, handler)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		filename := snaker.CamelToSnake(def.Name)
		filename = "rest_" + filename + ".go"

		f, err := os.Create(path + "/" + filename)
		if err != nil {
			return errors.Stack(err)
		}

		_, err = io.Copy(f, b)
		if err != nil {
			return errors.Stack(err)
		}

		err = f.Close()
		if err != nil {
			return errors.Stack(err)
		}

	}

	sort.Strings(registers)

	f, err := os.Create(path + "/main.go")
	if err != nil {
		return errors.Stack(err)
	}

	f.WriteString(rester.Main(registers))

	err = f.Close()
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
