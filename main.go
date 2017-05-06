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

var (
	pkgPath string
)

func main() {
	var (
		driver     *string = flag.String("driver", "mysql", "The SQL driver relevant to your request; one of 'lite', 'my', 'pg', or 'ms'.")
		pkgPathIn  *string = flag.String("pkg", "", "The package that you want to use for source material.")
		database   *string = flag.String("db", "", "The name of the database you want to analyze.")
		hostname   *string = flag.String("host", "", "The host (IP address or hostname) that hosts the database you want to analyze.")
		outputPath *string = flag.String("out", "", "The path where resulting files will be deposited.")
		password   *string = flag.String("pass", "", "The password of the user specificed by 'username'.")
		username   *string = flag.String("user", "", "The username on the database you want to analyze.")
		src        *string = flag.String("src", "", "The source of data that interests you; one of 'pkg' or 'db'.")
		dst        *string = flag.String("dst", "", "The desired output, one of 'types', 'control', 'bindings', 'rest', 'schema', or 'everything'.")
	)

	flag.Parse()

	switch *src {
	case "db":
		switch {
		case *username == "":
			flag.Usage()
			log.Println("When the 'src' argument is 'db', the 'user' argument is required.")
			os.Exit(1)
		case *password == "":
			flag.Usage()
			log.Println("When the 'src' argument is 'db', the 'pass' argument is required.")
			os.Exit(1)
		case *database == "":
			flag.Usage()
			log.Println("When the 'src' argument is 'db', the 'db' argument is required.")
			os.Exit(1)
		case *hostname == "":
			flag.Usage()
			log.Println("When the 'src' argument is 'db', the 'host' argument is required.")
			os.Exit(1)
		case *driver == "":
			flag.Usage()
			log.Println("When the 'src' argument is 'db', the 'driver' argument is required.")
			os.Exit(1)
		}
	case "pkg":
		switch {
		case *pkgPathIn == "":
			flag.Usage()
			log.Println("When the 'src' argument is 'pkg', the 'pkg' argument is required.")
			os.Exit(1)
		}
	default:
		flag.Usage()
		log.Println("The 'src' argument is required and must be one of 'pkg' or 'db'.")
		os.Exit(1)
	}

	// 'types', 'control', 'bindings', 'rest', or 'schema'
	switch *dst {
	case "types":
	case "control":
	case "bindings":
		switch {
		case *driver == "":
			flag.Usage()
			log.Println("When the 'dst' argument is 'everything', the 'schema' argument is required.")
			os.Exit(1)
		}
	case "http":
	case "schema":
		switch {
		case *driver == "":
			flag.Usage()
			log.Println("When the 'dst' argument is 'everything', the 'schema' argument is required.")
			os.Exit(1)
		}
	case "everything":
		switch {
		case *driver == "":
			flag.Usage()
			log.Println("When the 'dst' argument is 'everything', the 'driver' argument is required.")
			os.Exit(1)
		}
	default:
		flag.Usage()
		log.Println("The 'dst' argument is required and must be one of 'types',")
		log.Println("'control', 'bindings', 'rest', 'schema', or 'everything'.")
		os.Exit(1)
	}

	if *outputPath == "" {
		flag.Usage()
		log.Println("You must always specify an output directory.")
		log.Println("If you want to use the current directory, use '.' or './'.")
		os.Exit(1)
	}




	pkgs, err := extractor.ExtractPackage(*pkgPathIn)
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
				fmt.Println(pg.SelectOne(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectOneTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectMany(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(pg.SelectManyTx(pkg.Name, sdef))
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
				fmt.Println(mysql.SelectOne(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectOneTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectMany(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.SelectManyTx(pkg.Name, sdef))
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
				fmt.Println(mysql.UpdateMany(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(mysql.UpdateManyTx(pkg.Name, sdef))
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
				fmt.Println(sqlite.SelectOne(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectOneTx(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectMany(pkg.Name, sdef))
				fmt.Println()
				fmt.Println("/*============================================================================*/")
				fmt.Println()
				fmt.Println(sqlite.SelectManyTx(pkg.Name, sdef))
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
