package types

// Basic types
type BasicType struct {
	Name string
}

func (b BasicType) String() string { return b.Name }
func (b BasicType) Equals(other Type) bool {
	if bt, ok := other.(BasicType); ok {
		return b.Name == bt.Name
	}
	return false
}

// Predefined types
var (
	NumberType  = BasicType{"nambala"}
	BoolType    = BasicType{"mawu"}
	StringType  = BasicType{"mawu"}
	VoidType    = BasicType{"void"}
	AnyType     = BasicType{"any"}
	UnknownType = BasicType{"unknown"}
)
