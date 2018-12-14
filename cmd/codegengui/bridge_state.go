package main

import (
	"github.com/jackmanlabs/errors"
	"log"
)

func (b *QmlBridge) saveState() {
	err := saveState(b.PackageState())
	if err != nil {
		log.Print(errors.Stack(err))
	}
}

func (b *QmlBridge) updateState() {

	importPath := b.ImportPath()
	states := b.PackageState()

	state := PackageState{
		ImportPath:    importPath,
		BindingsPath:  b.BindingsPath(),
		SchemaPath:    b.SchemaPath(),
		InterfacePath: b.InterfacePath(),
		SqlDriver:     b.SqlDriver(),
		WriteTests:    b.IsWriteTests(),
		SelectedTypes: b.TypeTableModel().checkedTypes(),
	}

	// jlog.Log(state)

	states[importPath] = state

	b.SetPackageState(states)
}
