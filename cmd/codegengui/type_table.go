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

func (m *TypeTableModel) init() {
	m.modelData = []*TypeTableItem{}

	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectData(m.data)
	m.ConnectSetData(m.setData)
	m.ConnectFlags(m.flags)
}

func (m *TypeTableModel) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		RoleTypeName:        core.NewQByteArray2("TypeName", -1),
		RoleTypeDescription: core.NewQByteArray2("TypeDescription", -1),
		RoleTypeChecked:     core.NewQByteArray2("TypeSelected", -1),
	}
}

func (m *TypeTableModel) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *TypeTableModel) columnCount(*core.QModelIndex) int {
	return 2
}

func (m *TypeTableModel) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
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

func (m *TypeTableModel) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *TypeTableModel) add(name string, desc string, selected bool) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	item := &TypeTableItem{
		name:    name,
		desc:    desc,
		checked: selected,
	}
	m.modelData = append(m.modelData, item)
	m.EndInsertRows()
}

func (m *TypeTableModel) edit(name string, desc string, selected bool) {
	if len(m.modelData) == 0 {
		return
	}

	item := &TypeTableItem{
		name:    name,
		desc:    desc,
		checked: selected,
	}

	m.modelData[len(m.modelData)-1] = item

	m.DataChanged(
		m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()),
		m.Index(len(m.modelData)-1, 2, core.NewQModelIndex()),
		[]int{RoleTypeName, RoleTypeDescription, RoleTypeChecked})
}

func (m *TypeTableModel) setData(index *core.QModelIndex, variant *core.QVariant, role int) bool {

	item := m.modelData[index.Row()]
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

func (m *TypeTableModel) flags(index *core.QModelIndex) core.Qt__ItemFlag {
	return core.Qt__ItemIsEditable
}

func (m *TypeTableModel) checkRow(row int, checked bool) {
	item := m.modelData[row]
	item.checked = checked

	m.DataChanged(
		m.Index(row, 0, core.NewQModelIndex()),
		m.Index(row, 2, core.NewQModelIndex()),
		[]int{RoleTypeChecked})
}

func (m *TypeTableModel) toggleRow(row int) {
	item := m.modelData[row]
	item.checked = !item.checked

	m.DataChanged(
		m.Index(row, 0, core.NewQModelIndex()),
		m.Index(row, 2, core.NewQModelIndex()),
		[]int{RoleTypeChecked})
}

func (m *TypeTableModel) selectAll() {

	m.BeginResetModel()
	for _, item := range m.modelData {
		item.checked = true
	}
	m.EndResetModel()
}

func (m *TypeTableModel) selectNone() {
	m.BeginResetModel()
	for _, item := range m.modelData {
		item.checked = false
	}
	m.EndResetModel()
}

func (m *TypeTableModel) checkedTypes() []string {

	var names []string = make([]string, 0)
	for _, item := range m.modelData {
		if item.checked {
			names = append(names, item.name)
		}
	}

	return names
}
