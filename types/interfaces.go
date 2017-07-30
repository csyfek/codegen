package types

type SqlGenerator interface {
	Baseline() string
	Schema(pkg *Package) string
	Delete(def *Type) string
	DeleteTx(def *Type) string
	InsertOne(pkgName string, def *Type) string
	InsertOneTx(pkgName string, def *Type) string
	SelectMany(pkgName string, def *Type) string
	SelectManyTx(pkgName string, def *Type) string
	SelectOne(pkgName string, def *Type) string
	SelectOneTx(pkgName string, def *Type) string
	UpdateMany(pkgName string, def *Type) string
	UpdateManyTx(pkgName string, def *Type) string
	UpdateOne(pkgName string, def *Type) string
	UpdateOneTx(pkgName string, def *Type) string
}

type Extractor interface {
	Extract() (*Package, error)
}
