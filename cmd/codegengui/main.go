package main

import (
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
	"github.com/therecipe/qt/quickcontrols2"
	"os"
)

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	gui.NewQGuiApplication(len(os.Args), os.Args)

	bridge := NewBridge()

	quickcontrols2.QQuickStyle_SetStyle("Material")
	engine := qml.NewQQmlApplicationEngine(nil)
	engine.RootContext().SetContextProperty("bridge", bridge)
	engine.Load(core.NewQUrl3("qrc:/qml/main.qml", 0))
	gui.QGuiApplication_Exec()
}
