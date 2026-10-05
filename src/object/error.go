package object

import "fmt"

const ERROR_OBJ = "ERROR"

type Error struct {
	Object
	Message string
	Line    int
	Column  int
	File    string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }

func (e *Error) String() string {
	if e.Line > 0 {
		if e.File != "" {
			return fmt.Sprintf("%s:%d:%d: error: %s", e.File, e.Line, e.Column, e.Message)
		}
		return fmt.Sprintf("%d:%d: error: %s", e.Line, e.Column, e.Message)
	}
	return "error: " + e.Message
}

func (i *Error) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func NewError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}
