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
		if member.Name == name {
			return true
		}
	}

	return false
}

type Member struct {
	Name   string // Go-friendly (CamelCase) name.
	Type   string // All types are normalized to best available Go types.
	Length int    // Length is preserved for DB-specific configurations.
}
