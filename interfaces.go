package main


type SqlGenerator interface {
	Baseline() string
	Create(parent string, children []Column) string
	Delete(parent string, children []Column) string
	DeleteTx(parent string, children []Column) string
	InsertOne(pkgName string, parent string, children []Column) string
	InsertOneTx(pkgName string, parent string, children []Column) string
	SelectMany(pkgName string, parent string, children []Column) string
	SelectManyTx(pkgName string, parent string, children []Column) string
	SelectOne(pkgName string, parent string, children []Column) string
	SelectOneTx(pkgName string, parent string, children []Column) string
	UpdateMany(pkgName string, parent string, children []Column) string
	UpdateManyTx(pkgName string, parent string, children []Column) string
	UpdateOne(pkgName string, parent string, children []Column) string
	UpdateOneTx(pkgName string, parent string, children []Column) string
}

type Extractor interface{
	strings() ([]string, error)
	Children(parent string) ([]Column, error)
}

/*
We're making an interface for this because the various databases represent
columns in different ways.
*/
type Column interface {
	GoType() string
	GoName() string
	SqlType() string
	SqlName() string
	IsPrimary() bool
}