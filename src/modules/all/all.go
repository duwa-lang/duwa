package all

import (
	"github.com/sevenreup/duwa/src/modules/console"
	"github.com/sevenreup/duwa/src/modules/math"
	"github.com/sevenreup/duwa/src/object"
)

var Modules = map[string]*object.LibraryModule{}
var Functions = map[string]*object.LibraryFunction{}

func Builtins() map[string]object.LibraryModule {
	result := map[string]object.LibraryModule{
		"Khonso": *console.Module(),
		"Masamu": *math.Module(),
	}

	return result
}

func init() {
	for k, v := range console.Builtins() {
		Functions[k] = v
	}
	modules := Builtins()
	for k, v := range modules {
		Modules[k] = &v
	}
}
