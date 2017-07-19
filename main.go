package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/jackmanlabs/codegen/mssql"
	"github.com/jackmanlabs/codegen/mysql"
	"github.com/jackmanlabs/codegen/pg"
	"github.com/jackmanlabs/codegen/pkgex"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/codegen/types"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"io"
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

		dirstat, err := dir.Stat()
		if err != nil {
			log.Fatal(errors.Stack(err))
		}

		if !dirstat.IsDir() {
			log.Fatal("The output path must be a directory.")
		}

		log.Print("OUTPUT PATH: \t", *outputRoot)
	}

	if *src != "pkg" && (*dst == "types" || *dst == "everything") {
		*importPathTypes = packagePath(*outputRoot + "/types")
	}

	if len(*importPathTypes) != 0 {
		log.Print("PACKAGE PATH:\t", *importPathTypes)
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
		extractor = pkgex.NewExtractor(*importPathTypes)
	}

	var pkg *types.Package
	pkg, err := extractor.Extract()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	structMap := make(map[string]*types.Type)

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
		generateBindings(*outputRoot, *importPathTypes, generator, pkg)
	}
}

func mergeStructs(dst, src *types.Type) {
	for _, srcMember := range src.Members {
		exists := dst.ContainsMember(srcMember.GoName)
		if !exists {
			dst.Members = append(dst.Members, srcMember)
		}
	}
}

func generateBindings(outputRoot, importPathTypes string, generator types.SqlGenerator, pkg *types.Package) error {

	var (
		f    io.WriteCloser
		path string
	)

	path = outputRoot + "/data"
	d, err := os.Open(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir|os.ModePerm)
		if err != nil {
			return errors.Stack(err)
		}
	} else if err != nil {
		return errors.Stack(err)
	} else {
		d.Close()
	}

	// Write baseline file.

	f, err = os.Create(path + "/db.go")
	if err != nil {
		return errors.Stack(err)
	}

	f.Write([]byte(generator.Baseline()))

	err = f.Close()
	if err != nil {
		return errors.Stack(err)
	}

	for _, def := range pkg.Types {

		b := bytes.NewBuffer(nil)

		fmt.Fprint(b, `
package data

import(
	"database/sql"
	"github.com/jackmanlabs/errors"
)

`)

		if importPathTypes != "" {
			fmt.Fprintln(b)
			fmt.Fprintln(b, "import (")
			fmt.Fprintln(b, "\""+importPathTypes+"\"")
			fmt.Fprintln(b, ")")
			fmt.Fprintln(b)
		}

		fmt.Fprintln(b)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b, "// TABLE: "+def.Table)
		fmt.Fprintln(b, "// TYPE:  "+def.Name)
		fmt.Fprintln(b, "//##############################################################################")
		fmt.Fprintln(b)

		fmt.Fprint(b, generator.SelectOne(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.SelectOneTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.SelectMany(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.SelectManyTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.InsertOne(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.InsertOneTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateOne(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateOneTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateMany(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.UpdateManyTx(pkg.Name, def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.Delete(def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")
		fmt.Fprint(b, generator.DeleteTx(def))
		fmt.Fprint(b, "\n\n/*============================================================================*/\n\n")

		filename := snaker.CamelToSnake(def.Name)
		filename = "bindings_" + filename + ".go"

		f, err := os.Create(path + "/" + filename)
		if err != nil {
			return errors.Stack(err)
		}

		_, err = io.Copy(f, b)
		if err != nil {
			return errors.Stack(err)
		}

		err = f.Close()
		if err != nil {
			return errors.Stack(err)
		}
	}

	return nil
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
