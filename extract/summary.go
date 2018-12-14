package extract

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"unicode"

	"github.com/jackmanlabs/errors"
)

func Summary(importPath string) (map[string]string, error) {
	// This needs to be really fast.
	// I don't want to deal with UI delay when calling this.

	extractor := &summaryExtractor{
		Types: map[string]string{},
	}

	buildPkg, err := build.Import(importPath, "", 0)
	if err != nil {
		return nil, errors.Stack(err)
	}

	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fset, buildPkg.Dir, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, errors.Stack(err)
	}

	extractor.Fset = fset

	for _, astPkg := range astPkgs {
		if strings.HasSuffix(astPkg.Name, "_test") {
			continue
		}
		ast.Walk(extractor, astPkg)
	}

	return extractor.Types, nil
}

type summaryExtractor struct {
	//pkgPath string
	Types map[string]string
	Fset  *token.FileSet // for debugging
}

func (x *summaryExtractor) Visit(node ast.Node) (w ast.Visitor) {

	if x.Fset == nil {
		log.Println("fset is nil.")
		return nil
	}

	switch t := node.(type) {

	case *ast.TypeSpec:

		name := t.Name.String()
		firstChar := []rune(name)[0]
		if unicode.IsLower(firstChar) {
			return nil
		}

		typ := resolveTypeExpression(t.Type)

		x.Types[name] = typ

	case *ast.Field:

		return nil

	case *ast.FuncDecl:
		// Ignore function declarations.
		return nil

	case *ast.ImportSpec:

		return nil

	case *ast.ValueSpec:
		// Ignore constant and variable declarations.
		return nil

	case nil:

	default:
		//log.Printf("unexpected type %T\n", t) // %T prints whatever type t has

	}

	return x
}
