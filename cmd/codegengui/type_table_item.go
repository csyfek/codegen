package main

import "strings"

type TypeTableItem struct {
	name string
	desc string
}

type TypeTableItemsByName []TypeTableItem

// Len is the number of elements in the collection.
func (this TypeTableItemsByName) Len() int {
	return len(this)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (this TypeTableItemsByName) Less(i, j int) bool {
	a := this[i].name
	b := this[j].name
	return strings.Compare(a, b) < 0
}

// Swap swaps the elements with indexes i and j.
func (this TypeTableItemsByName) Swap(i, j int) {
	this[i], this[j] = this[j], this[i]
}
