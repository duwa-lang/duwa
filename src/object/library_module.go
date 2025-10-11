package object

import "maps"

const LIBRARY_MODULE = "LIBRARY_MODULE"

type LibraryModule struct {
	Name       string
	Methods    map[string]*LibraryFunction
	Properties map[string]Object
}

func (libraryModule *LibraryModule) String() string {
	return libraryModule.Name
}

func (libraryModule *LibraryModule) Type() ObjectType {
	return LIBRARY_MODULE
}

func (libraryModule *LibraryModule) Method(method string, args []Object) (Object, bool) {
	return nil, false
}

func NewBuiltInLibraryModule(name string, contents map[string]*LibraryFunction) *LibraryModule {
	methods := map[string]*LibraryFunction{}
	maps.Copy(methods, contents)
	m := &LibraryModule{
		Name:       name,
		Methods:    methods,
		Properties: map[string]Object{},
	}
	return m
}

func NewBuiltInLibraryModuleWithProperties(name string, methods map[string]*LibraryFunction, properties map[string]Object) *LibraryModule {
	methodsCopy := map[string]*LibraryFunction{}
	maps.Copy(methodsCopy, methods)

	propertiesCopy := map[string]Object{}
	maps.Copy(propertiesCopy, properties)

	m := &LibraryModule{
		Name:       name,
		Methods:    methodsCopy,
		Properties: propertiesCopy,
	}
	return m
}
