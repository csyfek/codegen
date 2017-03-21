package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/extractor"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/pg"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

func main() {
	var (
		doGolang *bool   = flag.Bool("go", false, "Generate Go code.")
		doSql    *bool   = flag.Bool("sql", false, "Generate SQL code.")
		doSqlite *bool   = flag.Bool("sqlite", false, "Use the MySQL dialect (default).")
		doMysql  *bool   = flag.Bool("mysql", false, "Use the MySQL dialect (default).")
		doPg     *bool   = flag.Bool("pg", false, "Use the PostgreSQL dialect.")
		pkgPath  *string = flag.String("pkg", "", "The package that you want to use for source material.")
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

	sqlCount := 0
	if *doMysql {
		sqlCount++
	}
	if *doSqlite {
		sqlCount++
	}
	if *doPg {
		sqlCount++
	}

	if sqlCount > 1 {
		log.Println("You must choose only one SQL dialect.")
		flag.Usage()
		os.Exit(1)
	}

	if sqlCount == 0 {
		*doMysql = true
	}

	//_, err := extractor.PackageStructs(*pkgPath)
	pkgs, err := extractor.PackageStructs(*pkgPath)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

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

	if *doGolang {
		if *doPg {
			fmt.Println(pg.Baseline())
			fmt.Println()
		} else if *doMysql {
			fmt.Println(mysql.Baseline())
			fmt.Println()
		} else if *doSqlite {
			fmt.Println(sqlite.Baseline())
			fmt.Println()
		}
	}

	for _, pkg := range pkgs {
		for _, sdef := range pkg.Structs {
			if *doSql && *doMysql {
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

			if *doSql && *doSqlite {
				fmt.Println("-- -----------------------------------------------------------------------------")
				fmt.Println()
				fmt.Println(sqlite.Create(sdef))
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

			if *doGolang && *doMysql {
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
				fmt.Println(mysql.UpdatePlural(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.UpdatePluralTx(pkg.Name, sdef))
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

			if *doGolang && *doSqlite {
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectSingular(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectSingularTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectPlural(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectPluralTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.Update(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.UpdateTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.Insert(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.InsertTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.Delete(sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.DeleteTx(sdef))
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
