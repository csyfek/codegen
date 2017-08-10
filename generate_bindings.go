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

func generateBindings(outputRoot, importPathTypes string, generator common.SqlGenerator, pkg *common.Package) error {

	var (
		f    io.WriteCloser
		path string
	)

	path = outputRoot + "/data"
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

	// Write baseline file.

	f, err = os.Create(path + "/bindings.go")
	if err != nil {
		return errors.Stack(err)
	}

	f.Write([]byte(generator.Baseline()))

	err = f.Close()
	if err != nil {
		return errors.Stack(err)
	}

	for _, def := range pkg.Types {

		if def.UnderlyingType != "struct" {
			continue
		}

		b := bytes.NewBuffer(nil)

		fmt.Fprint(b, `
package data

import(
	"database/sql"
	"github.com/jackmanlabs/errors"
)

`)

		if importPathTypes != "" {
			fmt.Fprintf(b, "\nimport \""+importPathTypes+"\"\n\n")
		}

		fmt.Fprintln(b)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b, "// TABLE: "+def.Table)
		fmt.Fprintln(b, "// TYPE:  "+def.Name)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b)

		fmt.Fprint(b, generator.SelectOne(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.SelectOneTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.SelectMany(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.SelectManyTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.InsertOne(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.InsertOneTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateOne(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateOneTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateMany(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateManyTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.Delete(def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.DeleteTx(def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")

		filename := snaker.CamelToSnake(def.Name)
		filename = "bindings_" + filename + ".go"

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
