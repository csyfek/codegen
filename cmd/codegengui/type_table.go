package main

import (
	"github.com/therecipe/qt/core"
)

func init() { TypeTable_QmlRegisterType2("CustomQmlTypes", 1, 0, "TypeTable") }

type TypeTable struct {
	core.QAbstractTableModel
	_ func()                                  `constructor:"init"`
	_ func()                                  `signal:"remove,auto"`
	_ func(item []*core.QVariant)             `signal:"add,auto"`
	_ func(firstName string, lastName string) `signal:"edit,auto"`

	modelData []TypeTableItem
}

func (m *TypeTable) init() {
	m.modelData = []TypeTableItem{}

	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectData(m.data)
}

func (m *TypeTable) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		RoleTypeName:        core.NewQByteArray2("TypeName", -1),
		RoleTypeDescription: core.NewQByteArray2("TypeDescription", -1),
	}
}

func (m *TypeTable) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *TypeTable) columnCount(*core.QModelIndex) int {
	return 2
}

func (m *TypeTable) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
	switch role {
	case RoleTypeName:
		return core.NewQVariant14(item.name)
	case RoleTypeDescription:
		return core.NewQVariant14(item.desc)
	}
	return core.NewQVariant()
}

func (m *TypeTable) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *TypeTable) add(item []*core.QVariant) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(m.modelData, TypeTableItem{item[0].ToString(), item[1].ToString()})
	m.EndInsertRows()
}

func (m *TypeTable) edit(firstName string, lastName string) {
	if len(m.modelData) == 0 {
		return
	}
	m.modelData[len(m.modelData)-1] = TypeTableItem{firstName, lastName}
	m.DataChanged(m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()), m.Index(len(m.modelData)-1, 1, core.NewQModelIndex()), []int{RoleTypeName, RoleTypeDescription})
}
