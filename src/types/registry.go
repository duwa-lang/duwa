package types

import "fmt"

type TypeRegistry struct {
	types map[string]Type
}

func NewTypeRegistry() *TypeRegistry {
	registry := &TypeRegistry{
		types: make(map[string]Type),
	}

	// Register built-in types
	registry.RegisterType("float", NumberType)
	registry.RegisterType("bool", BoolType)
	registry.RegisterType("string", StringType)
	registry.RegisterType("void", VoidType)

	return registry
}

func (tr *TypeRegistry) RegisterType(name string, typ Type) error {
	if _, exists := tr.types[name]; exists {
		return fmt.Errorf("type '%s' already exists", name)
	}
	tr.types[name] = typ
	return nil
}

func (tr *TypeRegistry) GetType(name string) (Type, bool) {
	typ, exists := tr.types[name]
	return typ, exists
}

func (tr *TypeRegistry) TypeExists(name string) bool {
	_, exists := tr.types[name]
	return exists
}
