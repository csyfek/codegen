package main

import (
	"flag"
	"github.com/jackmanlabs/codegen/common"
	"github.com/jackmanlabs/codegen/mssql"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/pg"
	"github.com/jackmanlabs/codegen/pkger"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/errors"
	"github.com/segmentio/go-camelcase"
	"github.com/serenize/snaker"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var (
		driver          *string = flag.String("driver", "mysql", "The SQL driver relevant to your request; one of 'sqlite', 'mysql', 'pg', or 'mssql'.")
		importPathTypes *string = flag.String("pkg", "", "The package that you want to use for source material.")
		database        *string = flag.String("db", "", "The name of the database you want to analyze.")
		hostname        *string = flag.String("host", "", "The host (IP address or hostname) that hosts the database you want to analyze.")
		outputRoot      *string = flag.String("out", "", "The path where resulting files will be deposited. Without a specified path, stdout will be used if possible.")
		password        *string = flag.String("pass", "", "The password of the user specified by 'username'.")
		username        *string = flag.String("user", "", "The username on the database you want to analyze.")
		src             *string = flag.String("src", "", "The source of data that interests you; one of 'pkg' or 'db'.")
		dst             *string = flag.String("dst", "", "The desired output, one of 'types', 'control', 'bindings', 'rest', 'schema', or 'everything'.")
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
		case *importPathTypes == "":
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
		case *outputRoot == "":
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

	if *outputRoot == "" {
		flag.Usage()
		log.Println("An output directory must be specified.")
		log.Println("If you want to use the current directory, use '.' or './'.")
		os.Exit(1)
	} else {
		*outputRoot = strings.TrimSuffix(*outputRoot, "/")
		dir, err := os.Open(*outputRoot)
		if err != nil {
			log.Fatal(errors.Stack(err))
		}
		defer dir.Close()

		dirStat, err := dir.Stat()
		if err != nil {
			log.Fatal(errors.Stack(err))
		}

		if !dirStat.IsDir() {
			log.Fatal("The output path must be a directory.")
		}

		log.Print("OUTPUT PATH: \t", *outputRoot)
	}

	if *src != "pkg" && (*dst == "types" || *dst == "everything") {
		*importPathTypes = packagePath(*outputRoot + "/types")
	}

	// TODO: Allow the user to specify existing type, control, and data packages.

	var importPathControl string
	if importPathControl == "" && (*dst == "control" || *dst == "everything") {
		importPathControl = packagePath(*outputRoot + "/control")
	}

	var importPathData string
	if importPathData == "" && (*dst == "control" || *dst == "everything") {
		importPathData = packagePath(*outputRoot + "/data")
	}

	var importPathFilters string
	if importPathFilters == "" && (*dst == "types" || *dst == "everything") {
		importPathFilters = packagePath(*outputRoot + "/filters")
	}

	if len(*importPathTypes) != 0 {
		log.Print("PACKAGE PATH TYPES:  \t", *importPathTypes)
	}

	if len(importPathControl) != 0 {
		log.Print("PACKAGE PATH CONTROL:\t", importPathControl)
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

	var extractor common.Extractor

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
		extractor = pkger.NewExtractor(*importPathTypes)
	}

	var pkg *common.Package
	pkg, err := extractor.Extract()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	// Only the package extractor generates a package name, but it's needed if we generate bindings, for example.
	if pkg.Name == "" {
		pkg.Name = "types"
	}

	structMap := make(map[string]*common.Type)

	// Prepare to flatten the structs for output generation by making them addressable.
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

	// Depending on the source, the SQL names may not be set. Set them now to default values.
	if *dst == "bindings" || *dst == "everything" || *dst == "schema" {
		for _, def := range pkg.Types {

			if def.Table == "" {
				def.Table = snaker.CamelToSnake(def.Name)
			}

			for mid, member := range def.Members {
				if member.SqlName == "" {
					def.Members[mid].SqlName = snaker.CamelToSnake(member.GoName)
				}
			}

		}
	}

	// We also need to set the JSON names to something convention-friendly.
	if *dst == "rest" || *dst == "everything" {
		for _, def := range pkg.Types {
			for m, member := range def.Members {
				def.Members[m].JsonName = camelcase.Camelcase(member.GoName)
			}
		}
	}

	var generator common.SqlGenerator
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
		generateSchema(*outputRoot, generator, pkg)
	}

	if *dst == "bindings" || *dst == "everything" {
		generateBindings(*outputRoot, []string{*importPathTypes, importPathFilters}, generator, pkg)
	}

	if *dst == "types" || *dst == "everything" {
		generateTypes(*outputRoot, pkg)
		generateFilters(*outputRoot, pkg)
	}

	if *dst == "rest" || *dst == "everything" {
		generateRest(*outputRoot, []string{*importPathTypes, importPathControl, importPathFilters}, pkg)
	}

	if *dst == "control" || *dst == "everything" {
		generateControls(*outputRoot, []string{*importPathTypes, importPathData, importPathFilters}, pkg)
	}
}

func mergeStructs(dst, src *common.Type) {
	for _, srcMember := range src.Members {
		exists := dst.ContainsMember(srcMember.GoName)
		if !exists {
			dst.Members = append(dst.Members, srcMember)
		}
	}
}

func packagePath(path string) string {

	var err error

	path, err = filepath.Abs(path)
	if err != nil {
		log.Print(errors.Stack(err))
		return ""
	}

	gopath := os.Getenv("GOPATH")
	gopath, err = filepath.Abs(gopath)
	if err != nil {
		log.Print(errors.Stack(err))
		return ""
	}

	if !strings.HasPrefix(path, gopath) {
		return ""
	}

	pkgpath := strings.TrimPrefix(path, gopath+"/src/")

	return pkgpath
}
