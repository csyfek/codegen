import QtQuick 2.10
import QtQuick.Controls 1.4
import QtQuick.Controls 2.3
import QtQuick.Layouts 1.3
import CustomQmlTypes 1.0

ApplicationWindow {
    width: 800
    height: 600
    visible: true
    title: qsTr("Transometron")

    GridLayout {
        columns: 2
        anchors.fill: parent

        TreeView {
            id: treeview
            Layout.fillWidth: true
            Layout.fillHeight: true
            frameVisible: true

            onClicked: bridge.selectPackage(index)
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

            model: bridge.packageTree
        }

        ColumnLayout {

            // anchors.fill: parent
            TableView {
                id: tableview

                Layout.fillWidth: true
                Layout.fillHeight: true

                model: bridge.typeTable

                TableViewColumn {
                    title: "Type Name"
                    role: "TypeName"
                }

                TableViewColumn {
                    title: "Description"
                    role: "TypeDescription"
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
