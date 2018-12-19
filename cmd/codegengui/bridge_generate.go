package main

import (
	"fmt"
	"log"

	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/codegen/extract"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/codegen/util"
	"github.com/jackmanlabs/errors"
	"github.com/therecipe/qt/widgets"
)

func (b *QmlBridge) generateBindings() {
	log.Print("Attempting to generate bindings.")

	generator := pickGenerator(b.SqlDriver())
	if generator == nil {
		widgets.QMessageBox_Warning(
			nil,
			"Driver Not Implemented",
			fmt.Sprintf("The SQL driver '%s' is not yet implemented.", b.SqlDriver()),
			widgets.QMessageBox__Ok,
			widgets.QMessageBox__Ok,
		)
	}

	pkg, err := extract.Extract(b.ImportPath())
	if err != nil {
		logAndShow(errors.Stack(err))
		return
	}

	// limit the models generated to the ones that were selected.
	names := b.TypeTableModel().checkedTypes()
	log.Print("Qty of models to process:", len(names))

	models := make([]*codegen.Model, 0)
	for _, model := range pkg.Models {
		if util.SetContainsString(names, model.Name) {
			models = append(models, model)
		}
	}

	var (
		// modelsSourcePath   string = pkg.AbsPath
		modelsImportPath   string = b.ImportPath()
		modelsPkgName      string = pkg.Name
		bindingsSourcePath string = b.BindingsPath()
		bindingsImportPath string
		bindingsPkgName    string
	)

	bindingsImportPath = util.ImportPath(bindingsSourcePath)
	bindingsPkgName = util.PackageName(bindingsImportPath)

	log.Print("Models Package Path: ", modelsImportPath)
	log.Print("Models Package Name: ", modelsPkgName)
	log.Print("Bindings Package Path: ", bindingsImportPath)
	log.Print("Bindings Package Name: ", bindingsPkgName)

	err = codegen.WriteBindings(
		generator,
		models,
		bindingsSourcePath,
		bindingsPkgName,
		modelsImportPath,
		modelsPkgName,
	)
	if err != nil {
		logAndShow(errors.Stack(err))
		return
	}

	err = codegen.WriteBindingsTests(
		generator,
		models,
		bindingsSourcePath,
		bindingsImportPath,
		bindingsPkgName,
		modelsImportPath,
		modelsPkgName,
	)
	if err != nil {
		logAndShow(errors.Stack(err))
		return
	}
}

func (b *QmlBridge) generateInterface() {}

func (b *QmlBridge) generateSchema() {}

func pickGenerator(driver string) codegen.SqlGenerator {
	var generator codegen.SqlGenerator
	switch driver {
	case "sqlite":
		generator = sqlite.NewGenerator()
	case "mysql":
		generator = mysql.NewGenerator()
	}

	return generator
}
