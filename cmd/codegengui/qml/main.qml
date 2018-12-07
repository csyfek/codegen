import QtQuick 2.10
//Item
import QtQuick.Controls 1.4
//TableView
import QtQuick.Controls 2.3
//Button
import QtQuick.Layouts 1.3
//ColumnLayout
import CustomQmlTypes 1.0


//CustomTableModel
Item {
    width: 800
    height: 600

    GridLayout {
        columns: 2
        anchors.fill: parent

        TreeView {
            Layout.fillWidth: true
            Layout.fillHeight: true
            frameVisible: true

            TableViewColumn {
                role: "PackageName"
                title: "Package"
                width: 100
            }
            TableViewColumn {
                role: "PackagePath"
                title: "Path"
                width: 200
            }
            model: PackageTree {
            }
        }

        ColumnLayout {

            // anchors.fill: parent
            TableView {
                id: tableview

                Layout.fillWidth: true
                Layout.fillHeight: true

                model: StructTable {
                }

                TableViewColumn {
                    role: "FirstName"
                    title: role
                }

                TableViewColumn {
                    role: "LastName"
                    title: role
                }
            }

            Button {
                Layout.fillWidth: true

                text: "remove last item"
                onClicked: tableview.model.remove()
            }

            Button {
                Layout.fillWidth: true

                text: "add new item"
                onClicked: tableview.model.add(["john", "doe"])
            }

            Button {
                Layout.fillWidth: true

                text: "edit last item"
                onClicked: tableview.model.edit("bob", "omb")
            }
        }
    }
}
