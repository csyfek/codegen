package codegen

type SqlGenerator interface {
	BindingsBaseline(bindingsPkgName string) []byte
	BindingsBaselineTests(importPaths []string, bindingsPackageName string, modelPackageName string) ([]byte, error)
	Bindings(importPaths []string, bindingsPackageName string, modelPackageName string, def *Model) ([]byte, error)
	BindingsTests(importPaths []string, bindingsPackageName string, modelPackageName string, def *Model) ([]byte, error)
	Schema(pkg *Package) ([]byte, error)
}
