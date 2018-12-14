package main

import (
	"flag"
	"github.com/jackmanlabs/codegen/extract"
	"github.com/jackmanlabs/codegen/util"
	"log"
	"os"

	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
)

func main() {
	var (
		driver             *string = flag.String("driver", "mysql", "The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'.")
		modelsImportPath   *string = flag.String("pkg", "", "The package that you want to use for source material.")
		bindingsSourcePath *string = flag.String("dst", "", "The desired output path of the bindings.")
		//database   *string = flag.String("db", "", "The name of the database you want to analyze.")
		//hostname   *string = flag.String("host", "", "The host (IP address or hostname) that hosts the database you want to analyze.")
		//outputRoot *string = flag.String("out", "", "The path where resulting files will be deposited. Without a specified path, stdout will be used if possible.")
		//password   *string = flag.String("pass", "", "The password of the user specified by 'username'.")
		//username   *string = flag.String("user", "", "The username on the database you want to analyze.")
		//modelsImportPath        *string = flag.String("modelsImportPath", "", "The source of data that interests you; one of 'pkg' or 'db'.")
	)

	flag.Parse()

	if *modelsImportPath == "" {
		flag.Usage()
		log.Println("The 'pkg' argument is required.")
		os.Exit(1)
	}

	switch *driver {
	case "sqlite":
	case "mysql":
	default:
		flag.Usage()
		log.Println("The 'driver' argument is required.")
		os.Exit(1)
	}

	if *bindingsSourcePath == "" {
		flag.Usage()
		log.Println("The 'bindingsSourcePath' argument is required.")
		os.Exit(1)
	}

	pkg, err := extract.Extract(*modelsImportPath)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	// Generate SQL Names
	for _, def := range pkg.Models {

		if def.Table == "" {
			def.Table = snaker.CamelToSnake(def.Name)
		}

		for mid, member := range def.Members {
			if member.SqlName == "" {
				def.Members[mid].SqlName = snaker.CamelToSnake(member.GoName)
			}
		}
	}

	var generator codegen.SqlGenerator
	switch *driver {
	case "sqlite":
		generator = sqlite.NewGenerator()
	case "mysql":
		generator = mysql.NewGenerator()
	}

	err = util.CheckDir(*bindingsSourcePath)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	var (
		//modelsImportPath   string = b.ImportPath()
		//BindingsAbsPath    string = b.BindingsPath()
		//modelsSourcePath   string = pkg.AbsPath
		modelsPkgName      string = pkg.Name
		bindingsImportPath string
		bindingsPkgName    string
	)

	bindingsImportPath = util.ImportPath(*bindingsSourcePath)
	bindingsPkgName = util.PackageName(bindingsImportPath)

	err = codegen.WriteBindings(generator, pkg.Models, *bindingsSourcePath, bindingsPkgName, *modelsImportPath, modelsPkgName)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	err = codegen.WriteBindingsTests(generator, pkg.Models, *bindingsSourcePath, bindingsImportPath, bindingsPkgName, *modelsImportPath, modelsPkgName)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}
}
