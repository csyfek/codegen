package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/extractor"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/pg"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

func main() {
	var (
		doGolang *bool   = flag.Bool("go", false, "Generate Go code.")
		doMy     *bool   = flag.Bool("my", false, "Use the MySQL dialect (default).")
		doPg     *bool   = flag.Bool("pg", false, "Use the PostgreSQL dialect.")
		doSql    *bool   = flag.Bool("sql", false, "Generate SQL code.")
		pkgPath  *string = flag.String("pkg", "", "The file that you want to analyze.")
	)

	flag.Parse()

	if *pkgPath == "" {
		log.Println("You must specify a Go package to analyze for this tool to work.")
		flag.Usage()
		os.Exit(1)
	}

	if !*doSql && !*doGolang {
		log.Println("You must specify which language output you desire.")
		flag.Usage()
		os.Exit(1)
	}

	if *doMy && *doPg {
		log.Println("You must choose only one SQL dialect.")
		flag.Usage()
		os.Exit(1)
	}

	if !*doMy && !*doPg {
		*doMy = true
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

	if *doGolang{
		// This makes it so that the syntax checker kicks in.
		fmt.Println(`
package main

import (
	"database/sql"
	"encoding/json"
	"github.com/jackmanlabs/errors"
)

func db() (*sql.DB, error) {
	return &sql.DB{}, nil
}

`)
	}

	for _, pkg := range pkgs {
		for _, sdef := range pkg.Structs {
			if *doSql && *doMy {
				fmt.Println("-- -----------------------------------------------------------------------------")
				fmt.Println()
				fmt.Println(mysql.Create(sdef))
				fmt.Println()
				//fmt.Println("-- -----------------------------------------------------------------------------")
			}

			if *doSql && *doPg {
				fmt.Println("-- -----------------------------------------------------------------------------")
				fmt.Println()
				fmt.Println(pg.Create(sdef))
				fmt.Println()
				//fmt.Println("-- -----------------------------------------------------------------------------")
			}



			if *doGolang && *doPg {
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectSingular(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectSingularTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectPlural(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectPluralTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.Update(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.UpdateTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.Insert(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.InsertTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.Delete(sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.DeleteTx(sdef))
				fmt.Println()
				//fmt.Println("/*============================================================================*/")
			}

			if *doGolang && *doMy {
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectSingular(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectSingularTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectPlural(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectPluralTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.Update(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.UpdateTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.Insert(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.InsertTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.Delete(sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.DeleteTx(sdef))
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
