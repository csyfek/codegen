package codegen

type SqlGenerator interface {
	Baseline() string
	Schema(pkg *Package) string
	Delete(def *Model) string
	DeleteTx(def *Model) string
	InsertOne(pkgName string, def *Model) string
	InsertOneTx(pkgName string, def *Model) string
	SelectMany(pkgName string, def *Model) string
	SelectManyTx(pkgName string, def *Model) string
	SelectOne(pkgName string, def *Model) string
	SelectOneTx(pkgName string, def *Model) string
	UpdateMany(pkgName string, def *Model) string
	UpdateManyTx(pkgName string, def *Model) string
	UpdateOne(pkgName string, def *Model) string
	UpdateOneTx(pkgName string, def *Model) string
}
