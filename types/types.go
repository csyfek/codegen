package types

type Package struct {
	Types   []*Type
	Imports map[string][]string
	Name    string
	Path    string
}

type Type struct {
	Name            string
	Members         []Member
	EmbeddedStructs []string
	UnderlyingType  string
}

func NewType() *Type {
	return &Type{
		Members:         make([]Member, 0),
		EmbeddedStructs: make([]string, 0),
	}
}

func (this *Type) ContainsMember(name string) bool {
	for _, member := range this.Members {
		if member.GoName() == name {
			return true
		}
	}

	return false
}

/*
Members are derived from different sources. Depending on the source, the Go and
SQL names will be generated differently, and, possibly, functionally.
*/
type Member interface {
	GoType() string
	GoName() string
	SqlType() string
	SqlName() string
	IsPrimary() bool
}
