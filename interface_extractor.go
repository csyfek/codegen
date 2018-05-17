package codegen

type Extractor interface {
	Extract() (*Package, error)
}
