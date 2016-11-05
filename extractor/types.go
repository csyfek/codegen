package extractor

import (
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

type StructDefinition struct {
	Name            string
	Members         []StructMemberDefinition
	embeddedStructs []string
}

type StructMemberDefinition struct {
	Name string
	Type string
}

type PackageDefinition struct {
	fset         *token.FileSet
	Name         string
	Structs      []*StructDefinition
	Imports      []string
	pendingIdent string
}

func (this *PackageDefinition) Visit(node ast.Node) ast.Visitor {
	switch t := node.(type) {
	case *ast.Package:
		this.Name = t.Name
	case *ast.FuncDecl:
		return nil
	case *ast.StructType:
		// push a new struct onto the slice.
		if this.Structs == nil {
			this.Structs = make([]*StructDefinition, 0)
		}
		this.Structs = append(this.Structs, &StructDefinition{Name: this.pendingIdent})
	case *ast.ValueSpec:
		// Ignore constant and variable declarations.
		return nil
	case *ast.ImportSpec:
		if this.Imports == nil {
			this.Imports = make([]string, 0)
		}
		this.Imports = append(this.Imports, t.Path.Value)
		return nil
	case *ast.GenDecl:
		if t.Tok != token.TYPE {
			return nil
		}
		for _, spec := range t.Specs {
			typeSpec := spec.(*ast.TypeSpec)
			this.pendingIdent = typeSpec.Name.String()
		}
	case *ast.Field:
		sdef := this.Structs[len(this.Structs)-1]

		var member StructMemberDefinition
		if len(t.Names) == 0 {
			printer.Fprint(os.Stderr, this.fset, t)
			return nil
		}
		member.Name = t.Names[0].String()
		member.Type = getExprType(t.Type)

		sdef.Members = append(sdef.Members, member)
	default:
	}

	return this
}
