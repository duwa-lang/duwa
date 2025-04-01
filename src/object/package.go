package object

import (
	"fmt"
)

const PACKAGE_OBJ = "PACKAGE"

type Package struct {
	Object
	Name   string
	Values map[string]Object
}

func NewPackageFromEnvironment(name string, env *Environment) *Package {
	values := make(map[string]Object)

	for k, v := range env.All() {
		values[k] = v
	}

	return &Package{
		Name:   name,
		Values: values,
	}
}

func (p *Package) Type() ObjectType {
	return PACKAGE_OBJ
}

func (p *Package) String() string {
	return p.Name
}

func (p *Package) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func (p *Package) GetPackageFunction(method string) (*Function, error) {
	function, ok := p.Values[method]
	if !ok {
		return nil, fmt.Errorf("function not found %s in package %s", method, p.Name)
	}

	fn, ok := function.(*Function)
	if !ok {
		return nil, fmt.Errorf("not a function %s in package %s", method, p.Name)
	}

	return fn, nil
}

func (p *Package) Get(name string) (Object, bool) {
	obj, ok := p.Values[name]
	return obj, ok
}
