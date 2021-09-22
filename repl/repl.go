package repl

import (
	"bufio"
	"fmt"
	"github.com/moreal/monkey/evaluator"
	"github.com/moreal/monkey/lexer"
	"github.com/moreal/monkey/object"
	"github.com/moreal/monkey/parser"
	"io"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, err io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	for {
		if _, err := fmt.Fprintf(err, PROMPT); err != nil {
			panic(err)
		}

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		evaluated := evaluator.Eval(p.ParseProgram(), env)
		if evaluated != nil {
			if _, err := fmt.Fprintln(out, evaluated.Inspect()); err != nil {
				panic(err)
			}
		}
	}
}
