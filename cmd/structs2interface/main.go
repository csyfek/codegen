package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/codegen/pkger"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"log"
	"os"
)

func main() {
	var (
		//driver     *string = flag.String("driver", "mysql", "The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'.")
		src *string = flag.String("pkg", "", "The package that you want to use for source material.")
		//database   *string = flag.String("db", "", "The name of the database you want to analyze.")
		//hostname   *string = flag.String("host", "", "The host (IP address or hostname) that hosts the database you want to analyze.")
		//outputRoot *string = flag.String("out", "", "The path where resulting files will be deposited. Without a specified path, stdout will be used if possible.")
		//password   *string = flag.String("pass", "", "The password of the user specified by 'username'.")
		//username   *string = flag.String("user", "", "The username on the database you want to analyze.")
		//src        *string = flag.String("src", "", "The source of data that interests you; one of 'pkg' or 'db'.")
		dst *string = flag.String("dst", "", "The desired output path of the bindings.")
	)

	flag.Parse()

	if *src == "" {
		flag.Usage()
		log.Println("The 'pkg' argument is required.")
		os.Exit(1)
	}

	if *dst == "" {
		flag.Usage()
		log.Println("The 'dst' argument is required.")
		os.Exit(1)
	}

	extractor := pkger.NewExtractor(*src)

	var pkg *codegen.Package
	pkg, err := extractor.Extract()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	err = codegen.CheckDir(*dst)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	// Write baseline file.

	f, err := os.Create(*dst + "/ds.go")
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	interfaceImportPath := codegen.ImportPath(*dst)
	interfacePkgName := codegen.PackageName(interfaceImportPath)
	log.Print("Package Path: ", interfaceImportPath)
	log.Print("Package Name: ", interfacePkgName)

	modelsPkgName := codegen.PackageName(*src)
	log.Print("Package Path: ", *src)
	log.Print("Package Name: ", modelsPkgName)

	pkgInterface, err := codegen.PackageInterface(pkg.Models, interfacePkgName)
	f.Write([]byte(pkgInterface))

	err = f.Close()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	for _, def := range pkg.Models {

		if def.UnderlyingType != "struct" {
			continue
		}

		if len(def.Members) == 0 {
			continue
		}

		out, err := codegen.ModelInterface([]string{*src}, interfacePkgName, modelsPkgName, def)
		if err != nil {
			log.Fatal(errors.Stack(err))
		}

		filename := fmt.Sprintf("/ds_%s.go", snaker.CamelToSnake(def.Name))

		f, err := os.Create(*dst + filename)
		if err != nil {
			log.Fatal(errors.Stack(err))
		}

		f.Write([]byte(out))

		err = f.Close()
		if err != nil {
			log.Fatal(errors.Stack(err))
		}
	}

}
