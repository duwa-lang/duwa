package evaluator

import (
	"os"

	"github.com/sevenreup/duwa/src/ast"
	"github.com/sevenreup/duwa/src/lexer"
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/parser"
)

var searchPaths []string
var imported map[string]*object.Environment

func resolveFilePath(node *ast.ImportStatement) (string, error) {
	return "", nil
}

func isImported(path string) bool {
	return false
}

func evaluateImportStatement(node *ast.ImportStatement, env *object.Environment) object.Object {
	path, err := resolveFilePath(node)
	if err != nil {
		return object.NewError("%s", err.Error())
	}

	// TODO: Check if named imports need to be merged
	if isImported(path) {
		return nil
	}

	return evaluateFile(path, node, env)
}

func evaluateFile(filePath string, node *ast.ImportStatement, env *object.Environment) object.Object {
	source, err := os.ReadFile(filePath)
	if err != nil {
		return newError("%d:%d:%s: runtime error: %s", node.Token.Pos.Line, node.Token.Pos.Column, node.Token.File, err.Error())
	}
	lexer := lexer.New(source)
	parser := parser.New(lexer)
	file := parser.ParseFile()

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

	if node.Type == ast.DefaultImport {
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
