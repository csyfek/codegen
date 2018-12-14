package codegen

import (
	"fmt"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"os"
)

func WriteBindings(generator SqlGenerator, models []*Model, modelsImportPath, modelsPkgName, bindingsSourcePath, bindingsPkgName string) error {

	// Write baseline file.

	f, err := os.Create(bindingsSourcePath + "/bindings.go")
	if err != nil {
		return errors.Stack(err)
	}

	_, err = f.Write([]byte(generator.BindingsBaseline(bindingsPkgName)))
	if err != nil {
		return errors.Stack(err)
	}

	err = f.Close()
	if err != nil {
		return errors.Stack(err)
	}

	// The actual bindings files.
	for _, def := range models {

		if def.UnderlyingType != "struct" {
			continue
		}

		if len(def.Members) == 0 {
			continue
		}

		out, err := generator.Bindings([]string{modelsImportPath}, bindingsPkgName, modelsPkgName, def)
		if err != nil {
			return errors.Stack(err)
		}

		filename := fmt.Sprintf("/bindings_%s.go", snaker.CamelToSnake(def.Name))

		f, err := os.Create(bindingsSourcePath + filename)
		if err != nil {
			return errors.Stack(err)
		}

		_, err = f.Write([]byte(out))
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

func WriteBindingsTests(generator SqlGenerator, models []*Model, modelsImportPath, modelsPkgName, bindingsSourcePath, bindingsImportPath, bindingsPkgName string) error {

	// Write baseline file for tests.

	f, err := os.Create(bindingsSourcePath + "/bindings_test.go")
	if err != nil {
		return errors.Stack(err)
	}

	baseline, err := generator.BindingsBaselineTests([]string{bindingsImportPath}, bindingsPkgName, modelsPkgName)
	if err != nil {
		return errors.Stack(err)
	}

	_, err = f.Write([]byte(baseline))
	if err != nil {
		return errors.Stack(err)
	}

	err = f.Close()
	if err != nil {
		return errors.Stack(err)
	}

	// The files containing the basic tests for the bindings.
	for _, def := range models {

		if def.UnderlyingType != "struct" {
			continue
		}

		if len(def.Members) == 0 {
			continue
		}

		out, err := generator.BindingsTests([]string{modelsImportPath}, bindingsPkgName, modelsPkgName, def)
		if err != nil {
			return errors.Stack(err)
		}

		filename := fmt.Sprintf("%s/bindings_%s_test.go", bindingsSourcePath, snaker.CamelToSnake(def.Name))

		f, err := os.Create(filename)
		if err != nil {
			return errors.Stack(err)
		}

		_, err = f.Write([]byte(out))
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
