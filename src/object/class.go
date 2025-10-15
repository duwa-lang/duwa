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
	// Create a new environment for each instance
	// Each instance gets its own environment with methods accessible via the class
	instanceEnv := CopyEnvironmentDefaults(c.Env)

	// Copy ALL properties and methods from the class to the instance
	// This ensures each instance has its own environment without shared references
	classProps := c.Env.All()
	for name, value := range classProps {
		instanceEnv.SetLocal(name, value)
	}

	// Set outer to class env so inherited methods/properties can still be found if needed
	instanceEnv.outer = c.Env

	instance := &Instance{Class: c, Env: instanceEnv}

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
