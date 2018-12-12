package main

import (
	"strings"
)

type TypeTableItem struct {
	name    string
	desc    string
	checked bool

	//core.QAbstractItemModel
	//_ string `property:"name"`
	//_ string `property:"desc"`
	//_ bool   `property:"checked"`
}

type TypeTableItemsByName []*TypeTableItem

// Len is the number of elements in the collection.
func (items TypeTableItemsByName) Len() int {
	return len(items)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (items TypeTableItemsByName) Less(i, j int) bool {
	a := items[i].name
	b := items[j].name
	return strings.Compare(a, b) < 0
}

// Swap swaps the elements with indexes i and j.
func (items TypeTableItemsByName) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}
