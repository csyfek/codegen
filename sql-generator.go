package main

import "github.com/jackmanlabs/codegen/extractor"

type SqlGenerator interface {
	Baseline() string
	Create(def *extractor.StructDefinition) string
	Delete(def *extractor.StructDefinition) string
	DeleteTx(def *extractor.StructDefinition) string
	InsertSingular(pkgName string, def *extractor.StructDefinition) string
	InsertSingularTx(pkgName string, def *extractor.StructDefinition) string
	SelectPlural(pkgName string, def *extractor.StructDefinition) string
	SelectPluralTx(pkgName string, def *extractor.StructDefinition) string
	SelectSingular(pkgName string, def *extractor.StructDefinition) string
	SelectSingularTx(pkgName string, def *extractor.StructDefinition) string
	UpdatePlural(pkgName string, def *extractor.StructDefinition) string
	UpdatePluralTx(pkgName string, def *extractor.StructDefinition) string
	UpdateSingular(pkgName string, def *extractor.StructDefinition) string
	UpdateSingularTx(pkgName string, def *extractor.StructDefinition) string
}
