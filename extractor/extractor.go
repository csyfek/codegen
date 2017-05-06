package extractor

import (
	"github.com/jackmanlabs/errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"unicode"
)

func ExtractPackage(pkgPath string) (*Package, error) {

	buildPkg, err := build.Import(pkgPath, "", 0)
	if err != nil {
		return nil, errors.Stack(err)
	}

	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fset, buildPkg.Dir, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, errors.Stack(err)
	}

	pkg := &Package{Fset: fset}
	for _, astPkg := range astPkgs {
		if strings.HasSuffix(astPkg.Name, "_test") {
			continue
		}
		ast.Walk(pkg, astPkg)
	}

	pkg.Path = pkgPath
	return pkg, nil
}

func (this *Package) Visit(node ast.Node) (w ast.Visitor) {

	if this.Fset == nil {
		log.Println("fset is nil.")
		return nil
	}

	switch t := node.(type) {

	case *ast.Package:
		this.Name = t.Name

	case *ast.TypeSpec:

		newType := NewType()
		newType.Name = t.Name.String()
		newType.UnderlyingType = resolveTypeExpression(t.Type)
		this.Types = append(this.Types, newType)

	case *ast.Field:

		parent := this.Types[len(this.Types)-1]
		typ := resolveTypeExpression(t.Type)

		// Handle embedded structs.
		if len(t.Names) == 0 {
			parent.EmbeddedStructs = append(parent.EmbeddedStructs, typ)
			return nil
		}

		name := t.Names[0].String()

		// Ignore fields that are not exported.
		name_ := []rune(name)
		if unicode.IsLower(name_[0]) {
			return nil
		}

		member := Member{
			Type: typ,
			Name: name,
		}

		parent.Members = append(parent.Members, member)

		return nil

	case *ast.FuncDecl:
		// Ignore function declarations.
		return nil

	case *ast.ImportSpec:
		path := t.Path.Value
		//path = strings.Trim(path, "\"")
		name := t.Name.Name
		names, ok := this.Imports[path]
		if !ok {
			names = make([]string, 0)
		}
		if !sContains(names, t.Name.Name) {
			names = append(names, name)
		}
		this.Imports[path] = names

		return nil

	case *ast.ValueSpec:
		// Ignore constant and variable declarations.
		return nil

	case nil:

	default:
		log.Printf("unexpected type %T\n", t) // %T prints whatever type t has

	}

	return this
}
