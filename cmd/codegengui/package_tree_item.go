package main

import "github.com/therecipe/qt/core"

type PackageTreeItem struct {
	core.QObject

	_ func() `constructor:"init"`

	_packageName string
	_packagePath string
	_parentItem  *PackageTreeItem

	_childItems []*PackageTreeItem
}

func (i *PackageTreeItem) init() {
	i.ConnectDestroyPackageTreeItem(i.destroyTreeItem)
}

func (i *PackageTreeItem) destroyTreeItem() {
	for _, child := range i._childItems {
		child.DestroyPackageTreeItem()
	}
	i.DestroyPackageTreeItemDefault()
}

func (i *PackageTreeItem) initWith(packageName, packagePath string) *PackageTreeItem {
	i._packageName = packageName
	i._packagePath = packagePath
	return i
}

func (i *PackageTreeItem) appendChild(child *PackageTreeItem) {
	child._parentItem = i
	i._childItems = append(i._childItems, child)
}

func (i *PackageTreeItem) withChildren(children []*PackageTreeItem) *PackageTreeItem{
	for _, child := range children{
		i.appendChild(child)
	}
	return i
}


func (i *PackageTreeItem) child(row int) *PackageTreeItem {
	return i._childItems[row]
}

func (i *PackageTreeItem) childCount() int {
	return len(i._childItems)
}

func (i *PackageTreeItem) columnCount() int {
	return 2
}

func (i *PackageTreeItem) data(column int) string {
	switch column {
	case 0:
		return i._packageName
	case 1:
		return i._packagePath
	}
	return ""
}

func (i *PackageTreeItem) row() int {
	if i._parentItem != nil {
		for index, item := range i._parentItem._childItems {
			if item.Pointer() == i.Pointer() {
				return index
			}
		}
	}
	return 0
}

func (i *PackageTreeItem) parentItem() *PackageTreeItem {
	return i._parentItem
}
