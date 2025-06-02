package types

// Type alias
type AliasType struct {
	Name           string
	UnderlyingType Type
}

func (a AliasType) String() string {
	return a.Name
}

func (a AliasType) Equals(other Type) bool {
	if at, ok := other.(AliasType); ok {
		return a.Name == at.Name
	}
	return false
}