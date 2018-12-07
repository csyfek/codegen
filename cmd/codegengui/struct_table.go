package main

import (
	"github.com/therecipe/qt/core"
)

func init() { StructTable_QmlRegisterType2("CustomQmlTypes", 1, 0, "StructTable") }

type StructTable struct {
	core.QAbstractTableModel
	_ func()                                  `constructor:"init"`
	_ func()                                  `signal:"remove,auto"`
	_ func(item []*core.QVariant)             `signal:"add,auto"`
	_ func(firstName string, lastName string) `signal:"edit,auto"`

	modelData []StructTableItem
}

func (m *StructTable) init() {
	m.modelData = []StructTableItem{{"john", "doe"}, {"john", "bob"}}

	m.ConnectRoleNames(m.roleNames)
	m.ConnectRowCount(m.rowCount)
	m.ConnectColumnCount(m.columnCount)
	m.ConnectData(m.data)
}

func (m *StructTable) roleNames() map[int]*core.QByteArray {
	return map[int]*core.QByteArray{
		RoleFirstName: core.NewQByteArray2("FirstName", -1),
		RoleLastName:  core.NewQByteArray2("LastName", -1),
	}
}

func (m *StructTable) rowCount(*core.QModelIndex) int {
	return len(m.modelData)
}

func (m *StructTable) columnCount(*core.QModelIndex) int {
	return 2
}

func (m *StructTable) data(index *core.QModelIndex, role int) *core.QVariant {
	item := m.modelData[index.Row()]
	switch role {
	case RoleFirstName:
		return core.NewQVariant14(item.firstName)
	case RoleLastName:
		return core.NewQVariant14(item.lastName)
	}
	return core.NewQVariant()
}

func (m *StructTable) remove() {
	if len(m.modelData) == 0 {
		return
	}
	m.BeginRemoveRows(core.NewQModelIndex(), len(m.modelData)-1, len(m.modelData)-1)
	m.modelData = m.modelData[:len(m.modelData)-1]
	m.EndRemoveRows()
}

func (m *StructTable) add(item []*core.QVariant) {
	m.BeginInsertRows(core.NewQModelIndex(), len(m.modelData), len(m.modelData))
	m.modelData = append(m.modelData, StructTableItem{item[0].ToString(), item[1].ToString()})
	m.EndInsertRows()
}

func (m *StructTable) edit(firstName string, lastName string) {
	if len(m.modelData) == 0 {
		return
	}
	m.modelData[len(m.modelData)-1] = StructTableItem{firstName, lastName}
	m.DataChanged(m.Index(len(m.modelData)-1, 0, core.NewQModelIndex()), m.Index(len(m.modelData)-1, 1, core.NewQModelIndex()), []int{RoleFirstName, RoleLastName})
}
