package main

import "github.com/therecipe/qt/core"

const (
	RoleFirstName = int(core.Qt__UserRole) + 1<<iota
	RoleLastName
	RolePackageName
	RolePackagePath
)
