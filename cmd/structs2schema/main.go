package main

import (
	"flag"
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/codegen/pkger"
	"github.com/jackmanlabs/codegen/sqlite"
	"github.com/jackmanlabs/errors"
	"github.com/serenize/snaker"
	"log"
	"os"
	"path/filepath"
	"strings"
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

	// Generate SQL Names
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

	generator := sqlite.NewGenerator()

	err = checkDir(*dst)
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	f, err := os.Create(*dst + "/schema.sql")
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	bindingsImportPath := importPath(*dst)
	bindingsPkgName := packageName(bindingsImportPath)
	log.Print("Package Path: ", bindingsImportPath)
	log.Print("Package Name: ", bindingsPkgName)

	modelsPkgName := packageName(*src)
	log.Print("Package Path: ", *src)
	log.Print("Package Name: ", modelsPkgName)

	f.Write([]byte(generator.Schema(pkg)))

	err = f.Close()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}
}

func checkDir(path string) error {
	d, err := os.Open(path)
	if os.IsNotExist(err) {
		err = os.Mkdir(path, os.ModeDir|os.ModePerm)
		if err != nil {
			log.Print(path)
			return errors.Stack(err)
		}
	} else if err != nil {
		return errors.Stack(err)
	} else {
		err = d.Close()
		if err != nil {
			return errors.Stack(err)
		}
	}

	return nil
}

func importPath(path string) string {

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

func packageName(packagePath string) string {
	chunks := strings.Split(packagePath, "/")
	return chunks[len(chunks)-1]
}
