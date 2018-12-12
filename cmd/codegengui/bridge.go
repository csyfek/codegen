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

	_ *PackageTree    `property:"packageTreeModel"`
	_ *TypeTableModel `property:"typeTableModel"`

	// Package Tree methods
	_ func(index *core.QModelIndex) `slot:"selectPackage,auto"`

	// Type Table methods
	_ func(row int) `slot:"selectType,auto"`

	// Processing parameters
	_ string `property:"schemaPath"`
	_ string `property:"bindingsPath"`
	_ string `property:"interfacePath"`
	_ string `property:"sqlDriver"`
}

func (ptr *QmlBridge) init() {
}

func NewBridge() *QmlBridge {
	bridge := NewQmlBridge(nil)

	packageTree := NewPackageTree(nil)
	bridge.SetPackageTreeModel(packageTree)

	typeTable := NewTypeTableModel(nil)
	bridge.SetTypeTableModel(typeTable)

	return bridge
}

func (ptr *QmlBridge) selectPackage(index *core.QModelIndex) {
	pkgPath := index.Data(RolePackagePath).ToString()
	tbl := ptr.TypeTableModel()
	if pkgPath == "" {
		tbl.BeginResetModel()
		tbl.modelData = []*TypeTableItem{}
		tbl.EndResetModel()
	} else {
		log.Print("Package checked:", pkgPath)
		names, err := extract.Names(pkgPath)
		if err != nil {
			log.Print(errors.Stack(err))
		}

		newItems := make([]*TypeTableItem, 0)
		for name, desc := range names {
			//log.Printf("TYPE: %s\t(%s)", name, desc)
			item := &TypeTableItem{
				name:    name,
				desc:    desc,
				checked: false,
			}
			newItems = append(newItems, item)
		}

		sort.Sort(TypeTableItemsByName(newItems))

		tbl.BeginResetModel()
		tbl.modelData = newItems
		tbl.EndResetModel()
	}
}

func (ptr *QmlBridge) selectType(row int) {
	table := ptr.TypeTableModel()
	item := table.modelData[row]

	log.Printf("Selected: %s (%t)", item.name, item.checked)
}

//func (b *QmlBridge) checkType(row int, checked bool) {
//	table := b.TypeTableModel()
//	item := table.modelData[row]
//
//	item.checked = checked
//
//	table.DataChanged(
//		table.Index(row, 0, core.NewQModelIndex()),
//		table.Index(row, 0, core.NewQModelIndex()),
//		[]int{RoleTypeChecked},
//	)
//}
