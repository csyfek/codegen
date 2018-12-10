package main

import "github.com/therecipe/qt/core"

const (
	RoleTypeName = int(core.Qt__UserRole) + 1<<iota
	RoleTypeDescription
	RolePackageName
	RolePackagePath
)
