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

	"github.com/sevenreup/duwa/src/lexer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	object.RegisterEvaluator(evaluator.Eval)
	scanner := bufio.NewScanner(in)
	env := object.Default()
	log := slog.Default()
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New([]byte(line))
		p := parser.New(l)
		file := p.ParseFile()
		if len(p.Errors()) != 0 {
			utils.PrintParserErrors(log, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(file, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.String())
			io.WriteString(out, "\n")
		}
	}
}
