package pkger

import (
	"github.com/jackmanlabs/codegen"
	"github.com/jackmanlabs/errors"
	"github.com/jackmanlabs/tipo"
	"github.com/serenize/snaker"
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

	parents := make([]codegen.Parent, 0)
	for _, name := range names {
		m, err := tipo.FindModel(name, this.pkgPath)
		if err != nil {
			return nil, errors.Stack(err)
		}

		name := m.Description

		children := make([]codegen.Child, 0)
		for n, c := range m.Children {
			child := codegen.Child{
				Model:    *c,
				GoName:   n,
				SqlName:  snaker.CamelToSnake(n),
				JsonName: snaker.CamelToSnake(n),
				GoType:   c.UnderlyingType,

				// These types are set and used by the individual SQL generators:
				//Length
				//SqlType
				//SqlConstraint
			}
			children = append(children, child)
		}

		parent := codegen.Parent{
			Model:    *m,
			Name:     name,
			Table:    codegen.Plural(m.Description),
			Children: children,
		}

		parents = append(parents, parent)
	}
	return parents, nil
}
