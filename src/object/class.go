package object

import "github.com/duwa-lang/duwa/src/ast"

const CLASS_OBJ = "CLASS"

type Class struct {
	Object
	Name *ast.Identifier
	Env  *Environment
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }

func (c *Class) String() string {
	return "class " + c.Name.String()
}

func (c *Class) CreateInstance(method string, args []Object) Object {
	instance := &Instance{Class: c, Env: c.Env}

	if ok := c.Env.Has("constructor"); ok {
		result := instance.Call("constructor", args)

		if result != nil && result.Type() == ERROR_OBJ {
			return result
		}
	}

	return instance
}

func (i *Class) Method(method string, args []Object) (Object, bool) {
	return nil, false
}
