package object

import (
	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/token"
)

var evaluator func(node ast.Node, env *Environment) Object

type ObjectType string

type Object interface {
	HasMethods
	Type() ObjectType
	String() string
}

type MapKey struct {
	Type  ObjectType
	Value uint64
}

type Mappable interface {
	MapKey() MapKey
}

type HasMethods interface {
	Method(method string, args []Object) (Object, bool)
}

type GoFunction func(env *Environment, tok token.Token, args ...Object) Object
type GoProperty func(env *Environment, tok token.Token) Object

func RegisterEvaluator(e func(node ast.Node, env *Environment) Object) {
	evaluator = e
}

func IsError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
