package main

import (
	"fmt"
	"log"
	"sort"

	"github.com/jackmanlabs/codegen/extract"
	"github.com/jackmanlabs/codegen/util"
	"github.com/jackmanlabs/errors"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/widgets"
)

// QmlBridge is specifically being used where the problem at hand cannot be
// solved by methods attached to QML models. For example, when one model needs
// to affect internal state of the Go application or multiple other models.
type QmlBridge struct {
	core.QObject

	_ func() `constructor:"init"`

	_ *PackageTree            `property:"packageTreeModel"`
	_ *TypeTableModel         `property:"typeTableModel"`
	_ map[string]PackageState `property:"packageState"`

	// General purpose methods
	_ func() `slot:"saveState,auto"`
	_ func() `slot:"updateState,auto"`

	// PackageTree methods
	_ func(index *core.QModelIndex) `slot:"selectPackage,auto"`

	// TypeTable methods
	_ func(row int) `slot:"selectType,auto"`
	_ func()        `slot:"generateBindings,auto"`

	// Processing parameters
	_ string `property:"bindingsPath"`
	_ string `property:"importPath"`
	_ string `property:"interfacePath"`
	_ string `property:"schemaPath"`
	_ string `property:"sqlDriver"`
	_ bool   `property:"writeTests"`
}

func (b *QmlBridge) init() {
}

func NewBridge() *QmlBridge {
	bridge := NewQmlBridge(nil)

	packageTree := NewPackageTree(nil)
	bridge.SetPackageTreeModel(packageTree)

	typeTable := NewTypeTableModel(nil)
	bridge.SetTypeTableModel(typeTable)

	return bridge
}

func (b *QmlBridge) selectPackage(index *core.QModelIndex) {
	importPath := index.Data(RolePackagePath).ToString()
	tbl := b.TypeTableModel()
	if importPath == "" {
		tbl.BeginResetModel()
		tbl.modelData = []*TypeTableItem{}
		tbl.EndResetModel()
	} else {
		log.Print("Package checked:", importPath)
		b.SetImportPath(importPath)
		names, err := extract.Summary(importPath)
		if err != nil {
			log.Print(errors.Stack(err))
		}

		state := b.PackageState()[importPath]
		b.SetBindingsPath(state.BindingsPath)
		b.SetImportPath(state.ImportPath)
		b.SetInterfacePath(state.InterfacePath)
		b.SetSchemaPath(state.SchemaPath)
		b.SetWriteTests(state.WriteTests)
		// TODO: make the sql driver dropdown useful.
		// b.SetSqlDriver(state.SqlDriver)

		newItems := make([]*TypeTableItem, 0)
		for name, desc := range names {

			// This application needs to ignore everything but structs.
			if desc != "struct" {
				continue
			}

			// log.Printf("TYPE: %s\t(%s)", name, desc)
			item := &TypeTableItem{
				name: name,
				desc: desc,
			}

			if util.SetContainsString(state.SelectedTypes, name) {
				item.checked = true
			}

			newItems = append(newItems, item)
		}

		sort.Sort(TypeTableItemsByName(newItems))

		tbl.BeginResetModel()
		tbl.modelData = newItems
		tbl.EndResetModel()
	}
}

func (b *QmlBridge) selectType(row int) {
	table := b.TypeTableModel()
	item := table.modelData[row]

	log.Printf("Selected: %s (%t)", item.name, item.checked)
}

// func (b *QmlBridge) checkType(row int, checked bool) {
// 	table := b.TypeTableModel()
// 	item := table.modelData[row]
//
// 	item.checked = checked
//
// 	table.DataChanged(
// 		table.Index(row, 0, core.NewQModelIndex()),
// 		table.Index(row, 0, core.NewQModelIndex()),
// 		[]int{RoleTypeChecked},
// 	)
// }

func logAndShow(err error) {
	log.Print(err)
	widgets.QMessageBox_Warning(
		nil,
		"Unexpected Error",
		fmt.Sprintf("An unexpected error occurred. Please see the log for details."),
		widgets.QMessageBox__Ok,
		widgets.QMessageBox__Ok,
	)
}
