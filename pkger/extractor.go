package pkger

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"log"
	"strings"
	"unicode"
	"github.com/jackmanlabs/tipo"
)

type extractorType struct {
	pkgPath string
}

func NewExtractor(pkgPath string) *extractorType {
	return &extractorType{
		pkgPath: pkgPath,
	}
}

func (this *extractorType) Extract() ([]codegen.Parent, error) {

	names, err := tipo.FindPackageModels(this.pkgPath)
if err != nil {
	return nil, errors.Stack(err)
}

models := make([]codegen.Parent, 0)
for _, name := range names{
m, err := tipo.FindModel(name, this.pkgPath)
if err != nil{
	return nil, errors.Stack(err)
}

model := codegen.NewModel()

}


}

func (this *extractorType) Visit(node ast.Node) (w ast.Visitor) {

	if this.Fset == nil {
		log.Println("fset is nil.")
		return nil
	}

	switch t := node.(type) {

	case *ast.Package:
		this.Name = t.Name

	case *ast.TypeSpec:

		newType := codegen.NewModel()
		newType.Name = t.Name.String()
		newType.UnderlyingType = resolveTypeExpression(t.Type)

		//log.Printf("GoType: %s\tUnderlyingType: %s", newType.Name, newType.UnderlyingType)

		this.Models = append(this.Models, newType)

	case *ast.Field:

		parent := this.Models[len(this.Models)-1]
		typ := resolveTypeExpression(t.Type)

		// Handle embedded structs.
		if len(t.Names) == 0 {
			parent.EmbeddedModels = append(parent.EmbeddedModels, typ)
			return nil
		}

		name := t.Names[0].String()

		// Ignore fields that are not exported.
		name_ := []rune(name)
		if unicode.IsLower(name_[0]) {
			return nil
		}

		member := codegen.Child{
			GoType: typ,
			GoName: name,
		}

		parent.Members = append(parent.Members, member)

		return nil

	case *ast.FuncDecl:
		// Ignore function declarations.
		return nil

	case *ast.ImportSpec:
		//ast.Print(this.Fset, t)

		var (
			path string = t.Path.Value
			name string
		)

		// Store the alias name if possible, otherwise empty string.
		if t.Name != nil {
			name = t.Name.Name
		}

		names, ok := this.Imports[path]
		if !ok {
			names = make([]string, 0)
		}
		if !sContains(names, name) {
			names = append(names, name)
		}
		this.Imports[path] = names

		return nil

	case *ast.ValueSpec:
		// Ignore constant and variable declarations.
		return nil

	case nil:

	default:
		//log.Printf("unexpected type %T\n", t) // %T prints whatever type t has

	}

	return this
}
