package evaluator

import (
	"os"
	"path/filepath"

	"github.com/duwa-lang/duwa/src/ast"
	"github.com/duwa-lang/duwa/src/modules"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/parser"
)

var imported map[string]*object.Environment

func isImported(path string) bool {
	_, ok := imported[path]
	return ok
}

func evaluateImportStatement(node *ast.ImportStatement, env *object.Environment) object.Object {
	isStd, path, err := resolveFilePath(node, env)
	if err != nil {
		return object.NewError("%s", err.Error())
	}

	// TODO: Check if named imports need to be merged
	if isImported(path) {
		return nil
	}

	if isStd {
		return handleStdImport(path, node, env)
	}

	return evaluateFile(path, node, env)
}

func handleStdImport(filePath string, node *ast.ImportStatement, env *object.Environment) object.Object {
	module, ok := modules.ImportModule(filePath)
	if !ok {
		return newError("%d:%d:%s: runtime error: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, "Failed to import std moduler")
	}

	env.Set(filePath, module)

	return nil
}

func evaluateFile(filePath string, node *ast.ImportStatement, env *object.Environment) object.Object {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return newError("%d:%d:%s: runtime error: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, err.Error())
	}

	parser := parser.NewParser()
	file := parser.ParseFile([]byte(source))

	if len(parser.Errors()) != 0 {
		for _, err := range parser.Errors() {
			env.Logger.Error(err)
		}
		return newError("%d:%d:%s: runtime error: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, parser.Errors()[0])
	}

	newEnvironment := object.CopyEnvironmentDefaults(env)

	result := Eval(file, newEnvironment)

	if isError(result) {
		return result
	}
	return addImportToEnvironment(node, env, newEnvironment)
}

func addImportToEnvironment(node *ast.ImportStatement, env *object.Environment, newEnvironment *object.Environment) object.Object {
	if node.Type == ast.DefaultImport {
		// Todo: fix this, instead pf have packages we should have modules
		newPackage := object.NewPackageFromEnvironment(node.DefaultAlias.Value, newEnvironment)
		env.Set(node.DefaultAlias.Value, newPackage)
	} else if node.Type == ast.NamedImport {
		for alias, export := range node.Exports {
			value, ok := newEnvironment.Get(export.Value)

			if !ok {
				return newError("%d:%d:%s: runtime error: %s is not exported", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, export.Value)
			}

			env.Set(alias, value)
		}
	} else {
		return newError("%d:%d:%s: runtime error: invalid import type %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, node.Type)
	}

	return nil
}

func resolveFilePath(node *ast.ImportStatement, env *object.Environment) (bool, string, error) {
	path := node.Module.Value
	if isStdImport(path) {
		return true, path, nil
	}
	if filepath.Ext(path) != ".duwa" {
		path += ".duwa"
	}
	if filepath.IsAbs(path) {
		return false, path, nil
	}

	filename, err := filepath.Abs(filepath.Join(env.GetDirectory(), path))
	if err != nil {
		return false, "", err
	}
	return false, filename, nil
}

func isStdImport(path string) bool {
	if filepath.Ext(path) == ".duwa" {
		path = path[:len(path)-len(".duwa")]
	}
	return modules.IsValidModuleImport(path)
}
