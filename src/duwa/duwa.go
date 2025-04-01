package duwa

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sevenreup/duwa/src/evaluator"
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/utils"

	"github.com/sevenreup/duwa/src/lexer"
	"github.com/sevenreup/duwa/src/parser"
)

type Duwa struct {
	Environment *object.Environment
}

func New(env *object.Environment) *Duwa {
	duwa := &Duwa{
		Environment: env,
	}
	duwa.registerEvaluator()
	return duwa
}

func (c *Duwa) RunFile(filePath string) object.Object {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	c.Environment.SetDirectory(filepath.Dir(filePath))
	return c.run(file)
}

func (c *Duwa) Run(data string) object.Object {
	return c.run([]byte(data))
}

func (c *Duwa) run(data []byte) object.Object {
	lexer := lexer.New(data)
	parser := parser.New(lexer)
	file := parser.ParseFile()
	if len(parser.Errors()) != 0 {
		utils.PrintParserErrors(c.Environment.Logger, parser.Errors())
		return nil
	}
	return evaluator.Eval(file, c.Environment)
}

func (c *Duwa) registerEvaluator() {
	object.RegisterEvaluator(evaluator.Eval)
}

func (c *Duwa) Call(function string, args []object.Object) object.Object {
	return c.Environment.Call(function, args)
}
