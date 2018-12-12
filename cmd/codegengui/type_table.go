package main

import (
	"github.com/therecipe/qt/core"
)

func init() { TypeTableModel_QmlRegisterType2("CustomQmlTypes", 1, 0, "TypeTableModel") }

type TypeTableModel struct {
	core.QAbstractTableModel
	_ func()                                                 `constructor:"init"`
	_ func()                                                 `slot:"remove,auto"`
	_ func()                                                 `slot:"selectAll,auto"`
	_ func()                                                 `slot:"selectNone,auto"`
	_ func(firstName string, lastName string, selected bool) `slot:"add,auto"`
	_ func(firstName string, lastName string, selected bool) `slot:"edit,auto"`
	_ func(row int, checked bool)                            `slot:"checkRow,auto"`
	_ func(row int)                                          `slot:"toggleRow,auto"`

	modelData []*TypeTableItem
}

func (ptr *TypeTableModel) init() {
	ptr.modelData = []*TypeTableItem{}

	ptr.ConnectRoleNames(ptr.roleNames)
	ptr.ConnectRowCount(ptr.rowCount)
	ptr.ConnectColumnCount(ptr.columnCount)
	ptr.ConnectData(ptr.data)
	ptr.ConnectSetData(ptr.setData)
	ptr.ConnectFlags(ptr.flags)
}

func (ptr *TypeTableModel) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		RoleTypeName:        core.NewQByteArray2("TypeName", -1),
		RoleTypeDescription: core.NewQByteArray2("TypeDescription", -1),
		RoleTypeChecked:     core.NewQByteArray2("TypeSelected", -1),
	}
}

func (ptr *TypeTableModel) rowCount(*core.QModelIndex) int {
	return len(ptr.modelData)
}

func (ptr *TypeTableModel) columnCount(*core.QModelIndex) int {
	return 2
}

func (ptr *TypeTableModel) data(index *core.QModelIndex, role int) *core.QVariant {
	item := ptr.modelData[index.Row()]
	switch role {
	case RoleTypeName:
		return core.NewQVariant14(item.name)
	case RoleTypeDescription:
		return core.NewQVariant14(item.desc)
	case RoleTypeChecked:
		return core.NewQVariant11(item.checked)
	}
	return core.NewQVariant()
}

func (ptr *TypeTableModel) remove() {
	if len(ptr.modelData) == 0 {
		return
	}
	ptr.BeginRemoveRows(core.NewQModelIndex(), len(ptr.modelData)-1, len(ptr.modelData)-1)
	ptr.modelData = ptr.modelData[:len(ptr.modelData)-1]
	ptr.EndRemoveRows()
}

func (ptr *TypeTableModel) add(name string, desc string, selected bool) {
	ptr.BeginInsertRows(core.NewQModelIndex(), len(ptr.modelData), len(ptr.modelData))
	item := &TypeTableItem{
		name:    name,
		desc:    desc,
		checked: selected,
	}
	ptr.modelData = append(ptr.modelData, item)
	ptr.EndInsertRows()
}

func (ptr *TypeTableModel) edit(name string, desc string, selected bool) {
	if len(ptr.modelData) == 0 {
		return
	}

	item := &TypeTableItem{
		name:    name,
		desc:    desc,
		checked: selected,
	}

	ptr.modelData[len(ptr.modelData)-1] = item

	ptr.DataChanged(
		ptr.Index(len(ptr.modelData)-1, 0, core.NewQModelIndex()),
		ptr.Index(len(ptr.modelData)-1, 2, core.NewQModelIndex()),
		[]int{RoleTypeName, RoleTypeDescription, RoleTypeChecked})
}

func (ptr *TypeTableModel) setData(index *core.QModelIndex, variant *core.QVariant, role int) bool {

	item := ptr.modelData[index.Row()]
	switch role {
	case RoleTypeName:
		v := variant.ToString()
		item.name = v
		return true
	case RoleTypeDescription:
		v := variant.ToString()
		item.desc = v
	case RoleTypeChecked:
		v := variant.ToBool()
		item.checked = v
	default:
		return false
	}
	return true
}

func (ptr *TypeTableModel) flags(index *core.QModelIndex) core.Qt__ItemFlag {
	return core.Qt__ItemIsEditable
}

func (ptr *TypeTableModel) checkRow(row int, checked bool) {
	item := ptr.modelData[row]
	item.checked = checked

	ptr.DataChanged(
		ptr.Index(row, 0, core.NewQModelIndex()),
		ptr.Index(row, 2, core.NewQModelIndex()),
		[]int{RoleTypeChecked})
}

func (ptr *TypeTableModel) toggleRow(row int) {
	item := ptr.modelData[row]
	item.checked = !item.checked

	ptr.DataChanged(
		ptr.Index(row, 0, core.NewQModelIndex()),
		ptr.Index(row, 2, core.NewQModelIndex()),
		[]int{RoleTypeChecked})
}

func (ptr *TypeTableModel) selectAll() {

	ptr.BeginResetModel()
	for _, item := range ptr.modelData {
		item.checked = true
	}
	ptr.EndResetModel()
}

func (ptr *TypeTableModel) selectNone() {
	ptr.BeginResetModel()
	for _, item := range ptr.modelData {
		item.checked = false
	}
	ptr.EndResetModel()
}
