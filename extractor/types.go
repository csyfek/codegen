package extractor

import (
	"fmt"
	"go/ast"
	"go/token"
)

type StructDefinition struct {
	Name            string
	Members         []StructMemberDefinition
	EmbeddedStructs []string
}

type StructMemberDefinition struct {
	Name string
	Type string
}

type PackageDefinition struct {
	Fset         *token.FileSet
	Name         string
	Structs      []*StructDefinition
	Imports      []string
	pendingIdent string
}

func (this *PackageDefinition) Visit(node ast.Node) ast.Visitor {

	if this.Fset == nil {
		fmt.Println("fset is nil.")
		return nil
	}

	switch t := node.(type) {
	case *ast.Package:
		this.Name = t.Name
	case *ast.FuncDecl:
		// Ignore function declarations.
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
		fmt.Print("*ast.Field\n")

		sdef := this.Structs[len(this.Structs)-1]

		// Handle embedded structs.
		if len(t.Names) == 0 {
			//ast.Fprint(os.Stdout, this.Fset, t, nil)

			if sdef.EmbeddedStructs == nil {
				sdef.EmbeddedStructs = make([]string, 0)
			}
			embeddedStruct := t.Type.(*ast.Ident).Name
			sdef.EmbeddedStructs = append(sdef.EmbeddedStructs, embeddedStruct)
			return nil
		}

		var member StructMemberDefinition
		member.Name = t.Names[0].String()
		member.Type = getExprType(t.Type)

		sdef.Members = append(sdef.Members, member)
		return nil
		//case nil:
		//default:
		//	fmt.Printf("unexpected type %T\n", t) // %T prints whatever type t has
	}

	return this
}
