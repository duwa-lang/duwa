package object

const INSTANCE_OBJ = "INSTANCE"

type Instance struct {
	Class *Class
	Env   *Environment
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }

func (i *Instance) String() string {
	return i.Class.Name.String()
}

func (i *Instance) Method(method string, args []Object) (Object, bool) {
	return nil, true
}

func (i *Instance) Call(method string, args []Object) Object {
	function, ok := i.Env.Get(method)
	if !ok {
		return NewError("undefined instance method call %s for %s", method, i.Class.Name.String())
	}
	methodFunction, ok := function.(*Function)
	if !ok {
		return NewError("Could not call instance method %s for %s", method, i.Class.Name.String())
	}

	methodEnv := createNewMethodInstanceEnvironment(methodFunction, args)
	return evaluator(methodFunction.Body, methodEnv)
}

func createNewMethodInstanceEnvironment(method *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(method.Env)

	for i, param := range method.Parameters {
		if len(args) > i {
			env.Set(param.Value, args[i])
		}
	}

	return env
}
