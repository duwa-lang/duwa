package repl

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"

	"github.com/sevenreup/duwa/src/evaluator"
	"github.com/sevenreup/duwa/src/object"
	"github.com/sevenreup/duwa/src/parser"
	"github.com/sevenreup/duwa/src/utils"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	object.RegisterEvaluator(evaluator.Eval)
	scanner := bufio.NewScanner(in)
	env := object.NewDefaultEnvironment()
	log := slog.Default()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		parser := parser.NewParser()
		file := parser.ParseFile([]byte(line))
		if len(parser.Errors()) != 0 {
			utils.PrintParserErrors(log, parser.Errors())
			continue
		}
		evaluated := evaluator.Eval(file, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.String())
			io.WriteString(out, "\n")
		}
	}
}
