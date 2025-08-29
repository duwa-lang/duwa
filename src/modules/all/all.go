package all

import (
	"github.com/sevenreup/duwa/src/modules/console"
	"github.com/sevenreup/duwa/src/modules/http"
	"github.com/sevenreup/duwa/src/modules/math"
	"github.com/sevenreup/duwa/src/object"
)

var Modules = map[string]*object.LibraryModule{}
var Functions = map[string]*object.LibraryFunction{}

func Builtins() map[string]object.LibraryModule {
	result := map[string]object.LibraryModule{
		"khonso": *console.Module(),
		"masamu": *math.Module(),
		"http":   *http.Module(),
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

func ImportModule(path string) (object.Object, bool) {
	module, ok := Modules[path]
	return module, ok
}

func IsValidModuleImport(path string) bool {
	_, ok := Modules[path]
	return ok
}
