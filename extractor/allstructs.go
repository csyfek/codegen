package extractor

import (
	"github.com/jackmanlabs/errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func PackageStructs(pkgPath string) ([]PackageDefinition, error) {

	var (
		bpkg *build.Package
		err  error
	)
	srcDirs := build.Default.SrcDirs()
	for _, srcDir := range srcDirs {
		bpkg, err = build.Import(pkgPath, srcDir, 0)
		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, errors.Stack(err)
	}

	pkgDefs, err := getStructs(bpkg.Dir)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return pkgDefs, nil
}

func getFolders(path string) ([]string, error) {
	//log.Printf("Descending into directory: %s\n", path)
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.Stack(err)
	}
	defer f.Close()

	dirs := []string{path}

	dirnames, err := f.Readdirnames(-1)
	for _, dirname := range dirnames {
		if dirname == "vendor" {
			continue
		}

		if strings.HasPrefix(dirname, ".") {
			continue
		}

		dirpath := path + "/" + dirname

		fi, err := os.Stat(dirpath)
		if err != nil {
			return nil, errors.Stack(err)
		}
		if !fi.IsDir() {
			continue
		}

		subdirs, err := getFolders(dirpath)
		if err != nil {
			return nil, errors.Stack(err)
		}

		for _, dir := range subdirs {
			dirs = append(dirs, dir)
		}
	}

	return dirs, nil
}

func getStructs(path string) ([]PackageDefinition, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, errors.Stack(err)
	}

	pkgDefs := make([]PackageDefinition, 0)
	for _, pkg := range pkgs {
		pkgDef := &PackageDefinition{fset: fset}
		ast.Walk(pkgDef, pkg)
		pkgDefs = append(pkgDefs, *pkgDef)
	}

	return pkgDefs, nil
}
