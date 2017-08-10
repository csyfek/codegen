package main

import (
	"github.com/jackmanlabs/codegen/common"
	"github.com/jackmanlabs/errors"
	"io"
	"os"
)

func generateSchema(outputRoot string, generator common.SqlGenerator, pkg *common.Package) error {

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

	f, err = os.Create(path + "/schema.sql")
	if err != nil {
		return errors.Stack(err)
	}

	f.Write([]byte(generator.Schema(pkg)))

	err = f.Close()
	if err != nil {
		return errors.Stack(err)
	}

	return nil
}
