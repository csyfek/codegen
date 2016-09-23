package structfinder

import (
	"fmt"
	"github.com/jackmanlabs/errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func NewStructFinderFromFile(filename string) (*StructFinder, error) {

	var (
		this *StructFinder = new(StructFinder)
		err  error
	)

	this.Data, err = ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Stack(err)
	}

	this.FSet = token.NewFileSet()
	this.File, err = parser.ParseFile(this.FSet, filename, nil, 0)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return this, nil
}

func NewStructFinderFromBytes(data []byte) (*StructFinder, error) {

	var (
		this *StructFinder = new(StructFinder)
		err  error
	)

	this.Data = data

	this.FSet = token.NewFileSet()
	this.File, err = parser.ParseFile(this.FSet, "", this.Data, 0)
	if err != nil {
		return nil, errors.Stack(err)
	}

	return this, nil
}

func (this *StructFinder) FindStructs() []StructDefinition {

	structDefs := make([]StructDefinition, 0)

	for _, dec := range this.File.Decls {

		switch decl := dec.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.TYPE {
				for _, spec_ := range decl.Specs {
					spec := spec_.(*ast.TypeSpec)

					structDef := StructDefinition{
						Members: make([]StructMemberDefinition, 0),
						Name:    spec.Name.String(),
						Package: this.File.Name.String(),
					}

					t := spec.Type.(*ast.StructType)

					for _, field := range t.Fields.List {
						structMemberDef := StructMemberDefinition{
							Name: field.Names[0].Name,
							Type: getExprType(field.Type),
						}
						structDef.Members = append(structDef.Members, structMemberDef)
					}

					structDefs = append(structDefs, structDef)
				}
			}
		}
	}

	return structDefs
}

type StructFinder struct {
	File *ast.File
	Data []byte
	FSet *token.FileSet
}

type StructDefinition struct {
	Package string
	Name    string
	Members []StructMemberDefinition
}

type StructMemberDefinition struct {
	Name string
	Type string
}

func getExprType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return "*" + getExprType(t.X)
	case *ast.ArrayType:
		return "[]" + getExprType(t.Elt)
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return getExprType(t.X) + "." + t.Sel.Name
	default:
		return fmt.Sprintf("Unknown Type: %T", t)
	}
}
