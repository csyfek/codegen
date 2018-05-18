package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/codegen/pkger"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"log"
	"os"
	"github.com/jackmanlabs/codegen/mysql"
)

func main() {
	var (
		driver     *string = flag.String("driver", "mysql", "The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'.")
		src *string = flag.String("pkg", "", "The package that you want to use for source material.")
	)

	flag.Parse()

	if *src == "" {
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


	extractor := pkger.NewExtractor(*src)

	var pkg *codegen.Package
	pkg, err := extractor.Extract()
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

	schema, err := generator.Schema(pkg)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	fmt.Print(schema)
}
