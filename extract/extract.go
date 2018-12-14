package extract

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/codegen/util"
	"github.com/jackmanlabs/errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"unicode"
)

type extractorType struct {
	pkgPath string
	codegen.Package
	fset *token.FileSet
}

func Extract(importPath string) (*codegen.Package, error) {
	ex := &extractorType{
		pkgPath: importPath,
	}

	buildPkg, err := build.Import(ex.pkgPath, "", 0)
	if err != nil {
		return nil, errors.Stack(err)
	}


	fset := token.NewFileSet()
	astPkgs, err := parser.ParseDir(fset, buildPkg.Dir, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return nil, errors.Stack(err)
	}

	ex.ImportPath = buildPkg.ImportPath
	ex.AbsPath = buildPkg.Dir
	ex.fset = fset
	ex.Imports = make(map[string][]string)

	for _, astPkg := range astPkgs {
		if strings.HasSuffix(astPkg.Name, "_test") {
			continue
		}
		ast.Walk(ex, astPkg)
	}

	return &ex.Package, nil
}

func (x *extractorType) Visit(node ast.Node) (w ast.Visitor) {

	if x.fset == nil {
		log.Println("fset is nil.")
		return nil
	}

	switch t := node.(type) {

	case *ast.Package:
		x.Name = t.Name

	case *ast.TypeSpec:

		newType := codegen.NewModel()
		newType.Name = t.Name.String()
		newType.UnderlyingType = resolveTypeExpression(t.Type)

		//log.Printf("GoType: %s\tUnderlyingType: %s", newType.Name, newType.UnderlyingType)

		x.Models = append(x.Models, newType)

	case *ast.Field:

		parent := x.Models[len(x.Models)-1]
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

		member := codegen.Member{
			GoType: typ,
			GoName: name,
		}

		parent.Members = append(parent.Members, member)

		return nil

	case *ast.FuncDecl:
		// Ignore function declarations.
		return nil

	case *ast.ImportSpec:
		//ast.Print(x.fset, t)

		var (
			path string = t.Path.Value
			name string
		)

		// Store the alias name if possible, otherwise empty string.
		if t.Name != nil {
			name = t.Name.Name
		}

		names, ok := x.Imports[path]
		if !ok {
			names = make([]string, 0)
		}
		if !util.SetContainsString(names, name) {
			names = append(names, name)
		}
		x.Imports[path] = names

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
