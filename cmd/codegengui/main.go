package main

import (
	"log"
	"os"

	"github.com/jackmanlabs/errors"
	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/qml"
	"github.com/therecipe/qt/quickcontrols2"
)

func main() {
	core.QCoreApplication_SetAttribute(core.Qt__AA_EnableHighDpiScaling, true)

	gui.NewQGuiApplication(len(os.Args), os.Args)

	// Load package states (if they exist)
	state, err := loadState()
	if err != nil {
		log.Fatal(errors.Stack(err))
	}

	bridge := NewBridge()
	bridge.SetPackageState(state)

	quickcontrols2.QQuickStyle_SetStyle("Material")
	engine := qml.NewQQmlApplicationEngine(nil)
	engine.RootContext().SetContextProperty("bridge", bridge)
	engine.Load(core.NewQUrl3("qrc:/qml/main.qml", 0))
	gui.QGuiApplication_Exec()
}
