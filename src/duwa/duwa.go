package duwa

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/duwa-lang/duwa/src/evaluator"
	"github.com/duwa-lang/duwa/src/object"
	"github.com/duwa-lang/duwa/src/utils"

	"github.com/duwa-lang/duwa/src/parser"
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
	parser := parser.NewParser()
	file := parser.ParseFile(data)
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


func (c *Duwa) CallE(function string, args []object.Object) (object.Object, error) {
	result := c.Environment.Call(function, args)
	if errObj, ok := result.(*object.Error); ok {
		return nil, errors.New(errObj.Message)
	}
	return result, nil
}
