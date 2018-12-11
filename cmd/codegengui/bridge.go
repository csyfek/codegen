package main

import (
	"log"
	"sort"

	"github.com/jackmanlabs/codegen/extract"
	"github.com/jackmanlabs/errors"
	"github.com/therecipe/qt/core"
)

type QmlBridge struct {
	core.QObject

	_ func() `constructor:"init"`

	_ *PackageTree `property:"packageTree"`
	_ *TypeTable   `property:"typeTable"`

	// Package Tree methods
	_ func(index *core.QModelIndex) `slot:"selectPackage,auto"`

	// Type Table methods
}

func (b *QmlBridge) init() {
}

func NewBridge() *QmlBridge {
	bridge := NewQmlBridge(nil)

	packageTree := NewPackageTree(nil)
	bridge.SetPackageTree(packageTree)

	typeTable := NewTypeTable(nil)
	bridge.SetTypeTable(typeTable)

	return bridge
}

func (b *QmlBridge) selectPackage(index *core.QModelIndex) {
	pkgPath := index.Data(RolePackagePath).ToString()
	tbl := b.TypeTable()
	if pkgPath == "" {
		tbl.BeginResetModel()
		tbl.modelData = []TypeTableItem{}
		tbl.EndResetModel()
	} else {
		log.Print("Package selected:", pkgPath)
		names, err := extract.Names(pkgPath)
		if err != nil {
			log.Print(errors.Stack(err))
		}

		newItems := make([]TypeTableItem, 0)
		for name, desc := range names {
			//log.Printf("TYPE: %s\t(%s)", name, desc)
			itm := TypeTableItem{
				name: name,
				desc: desc,
			}
			newItems = append(newItems, itm)
		}

		sort.Sort(TypeTableItemsByName(newItems))

		tbl.BeginResetModel()
		tbl.modelData = newItems
		tbl.EndResetModel()
	}
}
