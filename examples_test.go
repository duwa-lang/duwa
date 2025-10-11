package chewa

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/duwa-lang/duwa/src/evaluator"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/parser"
	"github.com/duwa-lang/duwa/src/utils"
	"github.com/duwa-lang/duwa/src/utils/environment"
	"github.com/shopspring/decimal"
)

func testEval(input string) object.Object {
	fileBytes, err := os.ReadFile(input)
	if err != nil {
		log.Fatal(err)
	}
	p := parser.NewParser()
	file := p.ParseFile(fileBytes)
	env := object.NewDefaultEnvironment()

	env.SetDirectory(filepath.Dir(input))

	evaluatorInstance := evaluator.Eval
	filename, _ := filepath.Abs("./")
	environment.SetCompilationSettings(filename)
	object.RegisterEvaluator(evaluatorInstance)

	return evaluator.Eval(file, env)
}

func TestImportExamples(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`./examples/imports.duwa`, 2},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			utils.TestIntegerObject(t, evaluated, decimal.NewFromInt(int64(integer)))
		} else {
			utils.TestNullObject(t, evaluated)
		}
	}
}
