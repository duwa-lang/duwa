package types

import "fmt"

// Array type
type ArrayType struct {
	ElementType Type
}

func (a ArrayType) String() string {
	return fmt.Sprintf("[]%s", a.ElementType.String())
}

func (a ArrayType) Equals(other Type) bool {
	if at, ok := other.(ArrayType); ok {
		return a.ElementType.Equals(at.ElementType)
	}
	return false
}
