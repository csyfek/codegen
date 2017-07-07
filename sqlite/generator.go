package sqlite

type generator bool

func NewGenerator() *generator {
	return new(generator)
}
