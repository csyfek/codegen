package mysql

type generator bool

func NewGenerator() *generator {
	return new(generator)
}
