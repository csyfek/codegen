package main

import (
	"bytes"
	"fmt"
	"github.com/jackmanlabs/codegen/common"
	"github.com/jackmanlabs/codegen/controler"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"io"
	"os"
)

func generateControls(outputRoot string, importPaths []string, pkg *common.Package) error {

	path := outputRoot + "/control"
	d, err := os.Open(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir|os.ModePerm)
		if err != nil {
			return errors.Stack(err)
		}
	} else if err != nil {
		return errors.Stack(err)
	} else {
		d.Close()
	}

	for _, def := range pkg.Types {
		if len(def.Members) == 0 {
			continue
		}

		b := bytes.NewBuffer(nil)

		fmt.Fprint(b, `
package control

import(
	"github.com/pborman/uuid"
	"github.com/jackmanlabs/errors"
`)

		for _, importPath := range importPaths {
			if importPath != "" {
				fmt.Fprintf(b, "\t\"%s\"\n", importPath)
			}
		}
		fmt.Fprint(b, ")\n\n")

		var control string

		fmt.Fprintln(b)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b, "// TYPE: "+def.Name)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b)

		control = controler.GetOne(def)
		fmt.Fprintln(b, control)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		control = controler.GetCollection(def)
		fmt.Fprintln(b, control)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		control = controler.Insert(def)
		fmt.Fprintln(b, control)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		control = controler.Upsert(def)
		fmt.Fprintln(b, control)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		control = controler.Delete(def)
		fmt.Fprintln(b, control)
		fmt.Fprint(b, "\n\n/*----------------------------------------------------------------------------*/\n\n")

		filename := snaker.CamelToSnake(def.Name)
		filename = "control_" + filename + ".go"

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

	return nil
}
