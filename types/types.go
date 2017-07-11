package types

type Package struct {
	Types   []*Type
	Imports map[string][]string
	Name    string
	Path    string
}

type Type struct {
	Name            string
	Table           string
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
		if member.GoName == name {
			return true
		}
	}

	return false
}

type Member struct {
	GoName  string // Go-friendly (CamelCase) name.
	SqlName string // Name to be used for SQL schemas and operations.
	Type    string // All types are normalized to best available Go types.
	Length  int    // Length is preserved for DB-specific configurations.

	// After much debate and trial-and-error, I've settled on keeping the GoName
	// and the SqlName distinct and accessible as strings. Therefore, the caller
	// should be responsible for knowing if and how these values are set and
	// used.
	//
	// The hope is that the Type and Member objects can be used for translations
	// both to and from SQL schemas/databases.
}
