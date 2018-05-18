package codegen

type SqlGenerator interface {
	BindingsBaseline(bindingsPkgName string) string
	BindingsBaselineTests(importPaths []string, bindingsPackageName string, modelPackageName string) (string, error)
	Bindings(importPaths []string, bindingsPackageName string, modelPackageName string, def *Model) (string, error)
	BindingsTests(importPaths []string, bindingsPackageName string, modelPackageName string, def *Model) (string, error)
	Schema(pkg *Package) (string, error)
}
