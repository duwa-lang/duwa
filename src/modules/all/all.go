package all

import (
	"github.com/duwa-lang/duwa/src/modules/console"
	"github.com/duwa-lang/duwa/src/modules/math"
	"github.com/duwa-lang/duwa/src/object"
	"maps"
)

var Modules = map[string]*object.LibraryModule{}
var Functions = map[string]*object.LibraryFunction{}

func Builtins() map[string]object.LibraryModule {
	result := map[string]object.LibraryModule{
		"khonso": *console.Module(),
		"masamu": *math.Module(),
	}

	return result
}

func init() {
	maps.Copy(Functions, console.Builtins())
	modules := Builtins()
	for k, v := range modules {
		Modules[k] = &v
	}
}

func RegisterModule(name string, module *object.LibraryModule) {
	Modules[name] = module
}

func RegisterFunction(name string, function *object.LibraryFunction) {
	Functions[name] = function
}

func ImportModule(path string) (object.Object, bool) {
	module, ok := Modules[path]
	return module, ok
}

func IsValidModuleImport(path string) bool {
	_, ok := Modules[path]
	return ok
}
