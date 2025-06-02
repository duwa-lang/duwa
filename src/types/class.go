package types

// Add custom types to your Type interface implementations
type CustomType struct {
	Name   string
	Fields map[string]Type // for struct-like types
}

func (c CustomType) String() string {
	return c.Name
}

func (c CustomType) Equals(other Type) bool {
	if ct, ok := other.(CustomType); ok {
		return c.Name == ct.Name
	}
	return false
}