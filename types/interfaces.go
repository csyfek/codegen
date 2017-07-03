package types

type SqlGenerator interface {
	Baseline() string
	Schema(typeName string, children []Member) string
	Delete(typeName string, children []Member) string
	DeleteTx(typeName string, children []Member) string
	InsertOne(pkgName string, typeName string, children []Member) string
	InsertOneTx(pkgName string, typeName string, children []Member) string
	SelectMany(pkgName string, typeName string, children []Member) string
	SelectManyTx(pkgName string, typeName string, children []Member) string
	SelectOne(pkgName string, typeName string, children []Member) string
	SelectOneTx(pkgName string, typeName string, children []Member) string
	UpdateMany(pkgName string, typeName string, children []Member) string
	UpdateManyTx(pkgName string, typeName string, children []Member) string
	UpdateOne(pkgName string, typeName string, children []Member) string
	UpdateOneTx(pkgName string, typeName string, children []Member) string
}

type Extractor interface {
	Extract() (*Package, error)
}
