package types

import (
	"fmt"
	"strings"
)

// Function type
type FunctionType struct {
	Params []Type
	Return Type
}

func (f FunctionType) String() string {
	params := make([]string, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.String()
	}
	return fmt.Sprintf("(%s) -> %s", strings.Join(params, ", "), f.Return.String())
}

func (f FunctionType) Equals(other Type) bool {
	if ft, ok := other.(FunctionType); ok {
		if len(f.Params) != len(ft.Params) {
			return false
		}
		for i, p := range f.Params {
			if !p.Equals(ft.Params[i]) {
				return false
			}
		}
		return f.Return.Equals(ft.Return)
	}
	return false
}
