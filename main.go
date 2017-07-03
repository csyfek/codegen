package main

import (
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/mssql"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/pg"
	"github.com/jackmanlabs/codegen/pkgex"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/codegen/types"
	"github.com/jackmanlabs/errors"
	"log"
	"os"
)

var (
	pkgPath string
)

func main() {
	var (
		driver     *string = flag.String("driver", "mysql", "The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'.")
		pkgPathIn  *string = flag.String("pkg", "", "The package that you want to use for source material.")
		database   *string = flag.String("db", "", "The name of the database you want to analyze.")
		hostname   *string = flag.String("host", "", "The host (IP address or hostname) that hosts the database you want to analyze.")
		outputPath *string = flag.String("out", "", "The path where resulting files will be deposited. Without a specified path, stdout will be used if possible.")
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
		case *outputPath == "" :
			flag.Usage()
			log.Println("Multiple languages can not be sanely printed to stdout.")
			log.Println("Please specify a target directory with the 'out' flag.")
			log.Println("If you want to use the current directory, use '.' or './'.")
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
		log.Println("Writing to stdout.")
		log.Println("If you want to dump to a directory, specify a target directory with the 'out' flag.")
		log.Println("If you want to use the current directory, use '.' or './'.")
	}

	if *driver != "" {
		switch *driver {
		case "sqlite":
		case "mysql":
		case "pg":
		case "mssql":
		default:
			flag.Usage()
			log.Println("The 'driver' argument must be one of 'sqlite', 'mysql', 'pg', or 'mssql'.")
			os.Exit(1)
		}
	}

	var pkg types.Package
	var extractor types.Extractor

	// No need to go crazy with validation; that's already been done above.
	if *src == "db" {
		switch *driver {
		case "sqlite":
			extractor = sqlite.NewExtractor(*database)
		case "mysql":
			extractor = mysql.NewExtractor(*username, *password, *hostname, *database)
		case "pg":
			extractor = pg.NewExtractor(*username, *password, *hostname, *database)
		case "mssql":
			extractor = mssql.NewExtractor(*username, *password, *hostname, *database)
		}
	} else {
		extractor = pkgex.NewExtractor(*pkgPathIn)
	}

	pkg, err := extractor.Extract()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	structMap := make(map[string]*types.Type)

	// Prepare to flatten the structs for MySQL generation by making them addressable.
	for _, s := range pkg.Types {
		structMap[s.Name] = s
	}

	// Merge the embedded structs into the parent structs.
	for _, parentStruct := range pkg.Types {
		for _, embeddedStructName := range parentStruct.EmbeddedStructs {
			embeddedStruct := structMap[embeddedStructName]
			mergeStructs(parentStruct, embeddedStruct)
		}
	}

	var generator types.SqlGenerator
	if *dst == "bindings" || *dst == "everything" || *dst == "schema" {
		switch *driver {
		case "sqlite":
			generator = sqlite.NewGenerator()
		case "mysql":
			generator = mysql.NewGenerator()
		case "pg":
			generator = pg.NewGenerator()
		case "mssql":
			generator = mssql.NewGenerator()
		}

	}

	if *dst == "schema" || *dst == "everything" {
		for _, sdef := range pkg.Types {

			fmt.Println("-- -----------------------------------------------------------------------------")
			fmt.Println()
			fmt.Println(generator.Schema(sdef))
			fmt.Println()
			//fmt.Println("-- -----------------------------------------------------------------------------")
		}
	}

	if *dst == "bindings" || *dst == "everything" {
		for _, sdef := range pkg.Types {

			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.SelectOne(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.SelectOneTx(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.SelectMany(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.SelectManyTx(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.UpdateOne(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.UpdateOneTx(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.UpdateMany(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.UpdateManyTx(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.InsertOne(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.InsertOneTx(pkg.Name, sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.Delete(sdef))
			fmt.Println()
			fmt.Println("/*============================================================================*/")
			fmt.Println()
			fmt.Println(generator.DeleteTx(sdef))
			fmt.Println()
		}
	}

}

func mergeStructs(dst, src *types.Type) {
	for _, srcMember := range src.Members {
		exists := dst.ContainsMember(srcMember.Name)
		if !exists {
			dst.Members = append(dst.Members, srcMember)
		}
	}
}
