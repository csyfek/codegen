import QtQuick 2.11
import QtQuick.Controls 1.4
import QtQuick.Controls 2.4
import QtQuick.Layouts 1.11
import QtQuick.Dialogs 1.3
import CustomQmlTypes 1.0

ApplicationWindow {
    width: 800
    height: 600
    visible: true
    title: qsTr("Transometron")

    TabView {
        anchors.fill: parent

        Tab {
            title: "Structs"

            GridLayout {
                columns: 2
                anchors.fill: parent

                TreeView {
                    id: packageTree
                    Layout.fillWidth: true
                    Layout.fillHeight: true
                    frameVisible: true

                    onClicked: {
                        bridge.selectPackage(index)
                        packageTree.resizeColumnsToContents()
                        typeTable.resizeColumnsToContents()
                    }

                    TableViewColumn {
                        id: colPackageName
                        role: "PackageName"
                        title: "Package"
                        width: 100
                    }
                    TableViewColumn {
                        id: colPackagePath
                        role: "PackagePath"
                        title: "Path"
                        width: 100
                    }

                    model: bridge.packageTreeModel
                }

                ColumnLayout {
                    RowLayout {
                        Button {
                            Layout.fillWidth: true
                            text: "Select All"
                            onClicked: bridge.typeTableModel.selectAll()
                        }
                        Button {
                            Layout.fillWidth: true
                            text: "Select None"
                            onClicked: bridge.typeTableModel.selectNone()
                        }
                    }

                    // anchors.fill: parent
                    TableView {
                        id: typeTable

                        Layout.fillWidth: true
                        Layout.fillHeight: true

                        // onClicked: bridge.typeTableModel.toggleRow(row)
                        model: bridge.typeTableModel

                        TableViewColumn {
                            title: ""
                            role: "TypeSelected"
                            width: 20

                            delegate: Rectangle {
                                anchors.fill: parent

                                // color: "blue"
                                CheckBox {
                                    scale: 0.50
                                    id: typeCheckBox
                                    anchors.centerIn: parent
                                    checked: styleData.value
                                    onClicked: {
                                        bridge.typeTableModel.checkRow(
                                                    styleData.row,
                                                    typeCheckBox.checked)
                                    }
                                }
                            }
                        }

                        TableViewColumn {
                            title: "Type Name"
                            role: "TypeName"
                        }

                        TableViewColumn {
                            title: "Description"
                            role: "TypeDescription"
                        }
                    }
                    ComboBox {
                        // currentIndex: 2
                        model: ["mysql", "sqlite", "mssql", "postgres"]
                        //width: 200
                        onCurrentIndexChanged: {
                            console.debug(model[currentIndex])
                            bridge.sqlDriver = model[currentIndex]
                        }
                    }
                    GridLayout {
                        columns: 3
                        width: parent.width
                        Button {
                            text: "Generate"
                            enabled: bridge.interfacePath !== ""
                        }

                        Button {
                            text: "Set"
                            onClicked: interfaceFileDialog.open()
                        }
                        Text {
                            textFormat: Text.RichText
                            text: "Interface Folder:<br>"
                                  + (bridge.interfacePath ? bridge.interfacePath : "not set")
                        }
                        Button {
                            text: "Generate"
                            enabled: bridge.bindingsPath !== ""
                        }
                        Button {
                            text: "Set"
                            onClicked: bindingsFileDialog.open()
                        }
                        Text {
                            textFormat: Text.RichText
                            text: "Binding Folder:<br>"
                                  + (bridge.bindingsPath ? bridge.bindingsPath : "not set")
                        }
                        Button {
                            text: "Generate"
                            enabled: bridge.schemaPath !== ""
                        }
                        Button {
                            text: "Set"
                            onClicked: schemaFileDialog.open()
                        }
                        Text {
                            textFormat: Text.RichText
                            text: "Schema File:<br>"
                                  + (bridge.schemaPath ? bridge.schemaPath : "not set")
                        }
                    }
                }
            }
        }
        Tab {
            title: "SQL DB"
        }
    }

    FileDialog {
        id: bindingsFileDialog
        title: "Choose the Go/SQL bindings directory"
        folder: shortcuts.home
        selectFolder: true
        selectExisting: false
        selectMultiple: false
        onAccepted: {
            bridge.bindingsPath = fileUrl
        }
    }
    FileDialog {
        id: schemaFileDialog
        title: "Choose the SQL schema file"
        folder: shortcuts.home
        selectExisting: false
        defaultSuffix: "sql"
        selectFolder: false
        selectMultiple: false
        onAccepted: {
            bridge.schemaPath = fileUrl
        }
    }
    FileDialog {
        id: interfaceFileDialog
        title: "Choose the Go data source interface directory"
        folder: shortcuts.home
        selectExisting: false
        selectFolder: true
        selectMultiple: false
        onAccepted: {
            bridge.interfacePath = fileUrl
        }
    }
}
