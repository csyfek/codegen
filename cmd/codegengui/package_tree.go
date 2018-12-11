package main

import (
	"log"

	"github.com/jackmanlabs/errors"
	"github.com/therecipe/qt/core"
)

func init() { PackageTree_QmlRegisterType2("CustomQmlTypes", 1, 0, "PackageTree") }

type PackageTree struct {
	core.QAbstractItemModel

	_ func() `constructor:"init"`

	_ func()                                       `signal:"remove,auto"`
	_ func(item []*core.QVariant)                  `signal:"add,auto"`
	_ func(packageName string, packagePath string) `signal:"edit,auto"`

	rootItem *PackageTreeItem
}

func (m *PackageTree) init() {
	m.rootItem = NewPackageTreeItem(nil)
	pkgs, err := findPackagePaths()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	for _, i := range pkgs {
		m.rootItem.appendChild(i)
	}

	m.ConnectIndex(m.index)
	m.ConnectParent(m.parent)
	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectData(m.data)
}

func (m *PackageTree) index(row int, column int, parent *core.QModelIndex) *core.QModelIndex {
	if !m.HasIndex(row, column, parent) {
		return core.NewQModelIndex()
	}

	var parentItem *PackageTreeItem
	if !parent.IsValid() {
		parentItem = m.rootItem
	} else {
		parentItem = NewPackageTreeItemFromPointer(parent.InternalPointer())
	}

	childItem := parentItem.child(row).Pointer()
	if childItem != nil {
		return m.CreateIndex(row, column, childItem)
	}
	return core.NewQModelIndex()
}

func (m *PackageTree) parent(index *core.QModelIndex) *core.QModelIndex {
	if !index.IsValid() {
		return core.NewQModelIndex()
	}

	item := NewPackageTreeItemFromPointer(index.InternalPointer())
	parentItem := item.parentItem()

	if parentItem.Pointer() == m.rootItem.Pointer() {
		return core.NewQModelIndex()
	}

	return m.CreateIndex(parentItem.row(), 0, parentItem.Pointer())
}

func (m *PackageTree) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		RolePackageName: core.NewQByteArray2("PackageName", -1),
		RolePackagePath: core.NewQByteArray2("PackagePath", -1),
	}
}

func (m *PackageTree) rowCount(parent *core.QModelIndex) int {
	if !parent.IsValid() {
		return m.rootItem.childCount()
	}
	return NewPackageTreeItemFromPointer(parent.InternalPointer()).childCount()
}

func (m *PackageTree) columnCount(parent *core.QModelIndex) int {
	if !parent.IsValid() {
		return m.rootItem.columnCount()
	}
	return NewPackageTreeItemFromPointer(parent.InternalPointer()).columnCount()
}

func (m *PackageTree) data(index *core.QModelIndex, role int) *core.QVariant {
	if !index.IsValid() {
		return core.NewQVariant()
	}

	item := NewPackageTreeItemFromPointer(index.InternalPointer())
	switch role {
	case RolePackageName:
		return core.NewQVariant14(item._packageName)
	case RolePackagePath:
		return core.NewQVariant14(item._packagePath)
	}
	return core.NewQVariant()
}

func (m *PackageTree) remove() {
	if m.rootItem.childCount() == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.rootItem._childItems)-1, len(m.rootItem._childItems)-1)
	item := m.rootItem._childItems[len(m.rootItem._childItems)-1]
	m.rootItem._childItems = m.rootItem._childItems[:len(m.rootItem._childItems)-1]
	m.EndRemoveRows()
	item.DestroyPackageTreeItem()
}

func (m *PackageTree) add(item []*core.QVariant) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.rootItem._childItems), len(m.rootItem._childItems))
	m.rootItem.appendChild(NewPackageTreeItem(nil).initWith(item[0].ToString(), item[1].ToString()))
	m.EndInsertRows()
}

func (m *PackageTree) edit(packageName string, packagePath string) {
	if m.rootItem.childCount() == 0 {
		return
	}
	m.BeginResetModel()
	item := m.rootItem._childItems[len(m.rootItem._childItems)-1]
	item._packageName = packageName
	item._packagePath = packagePath
	m.EndResetModel()

	//TODO:
	//ideally DataChanged should be used instead, but it doesn't seem to work ...
	//if you search for "qml treeview datachanged" online
	//it will just lead you to tons of unresolved issues
	//m.DataChanged(m.Index(item.row(), 0, core.NewQModelIndex()), m.Index(item.row(), 1, core.NewQModelIndex()), []int{FirstName, LastName})
	//feel free to send a PR, if you got it working somehow :)
}
