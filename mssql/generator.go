package mssql

type generator struct{}

func NewGenerator() *generator {
	return new(generator)
}
