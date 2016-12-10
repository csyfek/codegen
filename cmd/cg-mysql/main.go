package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/extractor"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

func main() {
	var (
		pkgPath  *string = flag.String("pkg", "", "The file that you want to analyze.")
		doSql    *bool   = flag.Bool("sql", false, "Generate SQL code.")
		doGolang *bool   = flag.Bool("go", false, "Generate Go code.")
	)

	flag.Parse()

	if *pkgPath == "" {
		log.Println("You must specify a Go file to analyze for this tool to work.")
		flag.Usage()
		os.Exit(1)
	}

	if !*doSql && !*doGolang {
		log.Println("You must specify which language output you desire.")
		flag.Usage()
		os.Exit(1)
	}

	//_, err := extractor.PackageStructs(*pkgPath)
	pkgs, err := extractor.PackageStructs(*pkgPath)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	//jlog.Log(pkgs)


	for _, pkg := range pkgs {

		structMap := make(map[string]*extractor.StructDefinition)

		// Prepare to flatten the structs for MySQL generation by making them addressable.
		for _, s := range pkg.Structs {
			structMap[s.Name] = s
		}

		// Merge the embedded structs into the parent structs.
		for _, parentStruct := range pkg.Structs {
			for _, embeddedStructName := range parentStruct.EmbeddedStructs {
				embeddedStruct := structMap[embeddedStructName]
				mergeStructs(parentStruct, embeddedStruct)
			}
		}
	}

	//jlog.Log(pkgs)

	for _, pkg := range pkgs {
		for _, sdef := range pkg.Structs {
			if *doSql {
				fmt.Println("-- -----------------------------------------------------------------------------")
				fmt.Println()
				fmt.Println(mysql.Create(sdef))
				fmt.Println()
				//fmt.Println("-- -----------------------------------------------------------------------------")
			}

			if *doGolang {
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectSingular(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectPlural(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.Update(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.Insert(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.Delete(sdef))
				fmt.Println()
				//fmt.Println("/*============================================================================*/")
			}
		}
	}
}

func mergeStructs(dst, src *extractor.StructDefinition) {
	for _, srcMember := range src.Members {
		exists := dst.ContainsMember(srcMember.Name)
		if !exists {
			dst.Members = append(dst.Members, srcMember)
		}
	}
}
