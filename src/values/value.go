package values

import "github.com/duwa-lang/duwa/src/object"

var (
	NULL     = &object.Null{}
	TRUE     = &object.Boolean{Value: true}
	FALSE    = &object.Boolean{Value: false}
	BREAK    = &object.Break{}
	CONTINUE = &object.Continue{}
)
